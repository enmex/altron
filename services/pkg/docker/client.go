package docker

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Client struct {
	client        *client.Client
	listenedLogs  map[string]io.ReadCloser
	listenedStats map[string]io.ReadCloser
	mut           *sync.Mutex
}

func NewClient() (*Client, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &Client{
		client: client,
		mut:    &sync.Mutex{},
	}, nil
}

func (c *Client) ContainerByPort(ctx context.Context, port uint16) (*types.Container, error) {
	containers, err := c.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	for _, c := range containers {
		for _, p := range c.Ports {
			if p.PrivatePort == port || p.PublicPort == port {
				return &c, nil
			}
		}
	}
	return nil, fmt.Errorf("container with port %d not found", port)
}

func (c *Client) ContainerByID(ctx context.Context, containerID string) (*types.Container, error) {
	containers, err := c.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	for _, c := range containers {
		if c.ID[:12] == containerID {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("container with id %s not found", containerID)
}

func (c *Client) ContainersPorts(ctx context.Context) (map[uint16]types.Container, error) {
	containers, err := c.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	containersPorts := make(map[uint16]types.Container)
	for _, c := range containers {
		for _, p := range c.Ports {
			containersPorts[p.PrivatePort] = c
			containersPorts[p.PublicPort] = c
		}
	}
	return containersPorts, nil
}

func (c *Client) Containers(ctx context.Context) ([]types.Container, error) {
	return c.client.ContainerList(ctx, types.ContainerListOptions{})
}

func (c *Client) Logs(ctx context.Context, containerID string, handler func(logData []byte) error) error {
	container, err := c.ContainerByID(ctx, containerID)
	if err != nil {
		return err
	}

	reader, err := c.client.ContainerLogs(ctx, container.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		return err
	}

	hdr := make([]byte, 8)
	for {
		_, err := reader.Read(hdr)
		if err == io.ErrClosedPipe || err == io.EOF {
			return handler([]byte("EOF"))
		}
		if err != nil {
			log.Fatal(err)
		}
		count := binary.BigEndian.Uint32(hdr[4:])
		dat := make([]byte, count)
		_, err = reader.Read(dat)
		if err == io.ErrClosedPipe || err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		bytes, err := json.Marshal(SendLogResponse{
			SentAt:  time.Now(),
			Message: base64.StdEncoding.EncodeToString(dat),
		})
		if err != nil {
			log.Fatal(err)
		}
		if err := handler(bytes); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) CloseLogs(containerID string) error {
	c.mut.Lock()
	defer c.mut.Unlock()
	_, ok := c.listenedLogs[containerID]
	if !ok {
		return nil
	}
	return c.listenedLogs[containerID].Close()
}

func (c *Client) ImagesByNetwork(ctx context.Context, networkName string) ([]types.ImageSummary, error) {
	return c.client.ImageList(ctx, types.ImageListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "network",
			Value: networkName,
		}),
	})
}

func (c *Client) ImageByName(ctx context.Context, imageName string) (*types.ImageSummary, error) {
	imageList, err := c.client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}
	for _, i := range imageList {
		for _, tag := range i.RepoTags {
			if tag == imageName {
				return &i, nil
			}
		}
	}
	return nil, fmt.Errorf("image with name %s not found", imageName)
}

func (c *Client) ContainersByNetwork(ctx context.Context, networkName string) ([]types.Container, error) {
	containers, err := c.client.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "network",
			Value: networkName,
		}),
	})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (c *Client) ContainerByImage(ctx context.Context, image string) (*types.Container, error) {
	var containers []types.Container
	var err error
	containers, err = c.Containers(ctx)
	if err != nil {
		return nil, err
	}
	for _, c := range containers {
		if c.Image == image {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("container with image %s not found", image)
}

func (c *Client) Health(ctx context.Context, containerName, containerID string) (<-chan ContainerHealth, error) {
	container, err := c.ContainerByID(ctx, containerID)
	if err != nil {
		return nil, err
	}
	stats, err := c.client.ContainerStats(ctx, containerID, true)
	if err != nil {
		return nil, err
	}

	reader := stats.Body
	healthChan := make(chan ContainerHealth)

	go func(reader io.ReadCloser) {
		decoder := json.NewDecoder(reader)
		lastCpuPercent := 0.0
		var lastMemoryUsage uint64 = 0
		var lastStats types.StatsJSON
		if err == io.ErrClosedPipe || err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		for {
			var stats types.StatsJSON
			err := decoder.Decode(&stats)
			if err == io.ErrClosedPipe || err == io.EOF {
				healthChan <- ContainerHealth{
					ID:     container.ID[:12],
					Name:   containerName,
					Status: "down",
				}
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			cpuPercent := c.calculateCPUPercentage(lastStats.CPUStats.CPUUsage.TotalUsage, lastStats.CPUStats.SystemUsage, &stats.CPUStats)
			networksStats := make(map[string]NetworkStats)
			for iface, networkStats := range stats.Networks {
				networksStats[iface] = NetworkStats{
					RxBytes:   networkStats.RxBytes,
					RxPackets: networkStats.RxPackets,
					TxBytes:   networkStats.TxBytes,
					TxPackets: networkStats.TxPackets,
				}
			}
			containerStats := Stats{
				PidsStats: stats.PidsStats,
				CPUUsage:  cpuPercent,
				MemoryStats: MemoryStats{
					Usage:    stats.MemoryStats.Usage,
					MaxUsage: stats.MemoryStats.Limit,
				},
				GrowthDynamics: GrowthDynamics{
					CPU:    c.compareValuesFloat64(cpuPercent, lastCpuPercent),
					Memory: c.compareValuesUint64(stats.MemoryStats.Usage, lastMemoryUsage),
				},
				Networks: networksStats,
			}
			var status string
			if cpuPercent >= 50.0 || float64(stats.MemoryStats.Usage)/float64(stats.MemoryStats.Limit) >= 0.5 {
				status = "mumble"
			} else {
				status = "up"
			}
			healthChan <- ContainerHealth{
				ID:     container.ID[:12],
				Name:   containerName,
				Status: status,
				Stats:  containerStats,
			}
			lastStats = stats
			lastCpuPercent = cpuPercent
			lastMemoryUsage = stats.MemoryStats.Usage
		}
	}(reader)
	return healthChan, nil
}

func (c *Client) CloseStats(containerID string) error {
	c.mut.Lock()
	defer c.mut.Unlock()
	_, ok := c.listenedStats[containerID]
	if !ok {
		return nil
	}
	return c.listenedStats[containerID].Close()
}

func (c *Client) calculateCPUPercentage(previousCPUUsage, previousSystemUsage uint64, stats *types.CPUStats) float64 {
	cpuPercent := 0.0

	cpuDelta := float64(stats.CPUUsage.TotalUsage) - float64(previousCPUUsage)
	systemDelta := float64(stats.SystemUsage) - float64(previousSystemUsage)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(stats.OnlineCPUs) * 100.0
	}
	return cpuPercent
}

func (c *Client) compareValuesUint64(a, b uint64) int8 {
	var diff int64 = int64(a - b)
	if diff == 0 {
		return 0
	}
	if diff < 0 {
		return -1
	}
	return 1
}

func (c *Client) compareValuesFloat64(a, b float64) int8 {
	diff := a - b
	if diff == 0 {
		return 0
	}
	if diff < 0 {
		return -1
	}
	return 1
}
