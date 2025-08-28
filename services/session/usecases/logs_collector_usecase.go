package usecases

import (
	"altron/pkg/amqp"
	"altron/pkg/sftp"
	"altron/session/generated/spec"
	"altron/session/interfaces"
	"altron/session/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var logsCacheSize = 50

type LogsData struct {
	ContainerID string
	Logs        []string
	Size        int64
}

type CachedLogs struct {
	Logs []string
	Mu   *sync.Mutex
}

var _ interfaces.LogsCollectorUseCase = (*LogsCollectorUseCase)(nil)

type LogsCollectorUseCase struct {
	log        *logrus.Logger
	ftpClient  *sftp.Client
	producer   *amqp.Producer
	cachedLogs sync.Map
}

func NewLogsCollectorUseCase(log *logrus.Logger, sftpClient *sftp.Client, client *amqp.Client) (*LogsCollectorUseCase, error) {
	producer, err := amqp.NewProducer(client)
	if err != nil {
		return nil, err
	}

	return &LogsCollectorUseCase{
		log:        log,
		ftpClient:  sftpClient,
		producer:   producer,
		cachedLogs: sync.Map{},
	}, nil
}

func (u *LogsCollectorUseCase) AddLogsContainer(containerID string) error {
	u.cachedLogs.LoadOrStore(containerID, &CachedLogs{
		Logs: make([]string, 0),
		Mu:   &sync.Mutex{},
	})
	return u.producer.CreateExchange(fmt.Sprintf("logs-%s", containerID))
}

func (u *LogsCollectorUseCase) DeleteLogsContainer(containerID string) error {
	u.cachedLogs.Delete(containerID)
	return u.producer.DeleteExchange(fmt.Sprintf("logs-%s", containerID))
}

func (u *LogsCollectorUseCase) MonitorLogsJson(ctx context.Context) error {
	for {
		time.Sleep(time.Second)
		list, err := u.ftpClient.List("logs")
		if err != nil {
			return err
		}
		if len(list) == 0 {
			continue
		}
		wg := &sync.WaitGroup{}
		wg.Add(len(list))
		for _, entry := range list {
			path := "/logs/" + entry
			if err := u.ftpClient.Download(path); err != nil {
				return fmt.Errorf("error downloading entry %s: %v", entry, err)
			}
			go u.processLogsJson(path, wg)
		}
		wg.Wait()
	}
}

func (u *LogsCollectorUseCase) processLogsJson(jsonPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer os.Remove(jsonPath)
	defer u.ftpClient.Delete(jsonPath)
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		u.log.Errorln(err)
		return
	}
	defer jsonFile.Close()

	data, err := io.ReadAll(jsonFile)
	if err != nil {
		u.log.Errorln(err)
		return
	}

	var logsData LogsData
	if err := json.Unmarshal(data, &logsData); err != nil {
		u.log.Errorln(err)
		return
	}
	v, ok := u.cachedLogs.Load(logsData.ContainerID)
	if !ok {
		return
	}
	cache := v.(*CachedLogs)
	cache.Mu.Lock()
	cache.Logs = append(cache.Logs, logsData.Logs...)
	if len(cache.Logs) > logsCacheSize {
		cache.Logs = cache.Logs[len(cache.Logs)-logsCacheSize:]
	}
	cache.Mu.Unlock()
	for _, log := range logsData.Logs {
		if err := u.producer.SendMessage(
			fmt.Sprintf("logs-%s", logsData.ContainerID),
			[]byte(log),
		); err != nil {
			u.log.Errorln(err)
			return
		}
	}
	u.log.Infof("logs data for container %s has been read and processed", logsData.ContainerID)
}

func (u *LogsCollectorUseCase) GetCachedLogs(containerID string) (*spec.GetLatestContainerLogsResponse, error) {
	v, ok := u.cachedLogs.Load(containerID)
	if !ok {
		return nil, models.ErrorCachedLogsNotFound
	}
	cache := v.(*CachedLogs)
	cache.Mu.Lock()
	defer cache.Mu.Unlock()

	return &spec.GetLatestContainerLogsResponse{
		Logs: cache.Logs,
	}, nil
}

func (u *LogsCollectorUseCase) UpdateContainer(ctx context.Context, request *spec.UpdateContainerRequest) error {
	if request.OldContainer != nil {
		if err := u.DeleteLogsContainer(*request.OldContainer); err != nil {
			return err
		}
	}
	return u.AddLogsContainer(request.NewContainer)
}
