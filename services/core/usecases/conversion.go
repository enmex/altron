package usecases

import (
	"altron/config"
	"altron/core/dto"
	"altron/core/generated/spec"
	"altron/core/interfaces"
	"altron/core/models"
	req "altron/pkg/request"
	"altron/pkg/sftp"
	"altron/utils"
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

var _ interfaces.ConversionUseCase = (*ConversionUseCase)(nil)

type ExploitTemplatePayload struct {
	AltronHost  string
	ServicePort uint16
	Exploit     string
}

type ConversionUseCase struct {
	log *logrus.Logger
	cfg        *config.AppConfig
	sftpClient *sftp.Client
}

func NewConversionUseCase(log *logrus.Logger, cfg *config.AppConfig, sftpClient *sftp.Client) *ConversionUseCase {
	return &ConversionUseCase{
		log:        log,
		cfg:        cfg,
		sftpClient: sftpClient,
	}
}

func (u *ConversionUseCase) ConvertSessionToExploit(ctx context.Context, request *spec.ConvertSessionToExploitRequest) (*spec.ConvertToExploitResponse, error) {
	if len(request.Session.Packets) > 50 {
		request.Session.Packets = append(request.Session.Packets[:40], request.Session.Packets[:len(request.Session.Packets)-10]...)
	}

	if strings.HasPrefix(request.ExportType, "pwntools") {
		return u.convertSessionToPwntools(ctx, request)
	}

	payloads := make([]string, 0, len(request.Session.Packets))
	for _, packet := range request.Session.Packets {
		if packet.IsRequest {
			payloads = append(payloads, packet.Payload)
		}
	}

	res, err := req.Post[dto.ConvertToExploitResponse](
		fmt.Sprintf("http://altron.converter.loc:%d/convert_to_requests_sploit", u.cfg.AltronConverterPort),
		dto.ConvertSessionToExploitRequest{
			Base64Strings: payloads,
		},
	)
	if err != nil {
		return nil, err
	}

	return &spec.ConvertToExploitResponse{
		Exploit: res.Result,
	}, nil
}

func (u *ConversionUseCase) ConvertPacketToExploit(ctx context.Context, request *spec.ConvertPacketToExploitRequest) (*spec.ConvertToExploitResponse, error) {
	if strings.EqualFold(request.ExportType, "ascii") {
		return u.convertToAsciiBytes(ctx, request)
	}
	if strings.EqualFold(request.ExportType, "plugin") {
		return u.convertToPlugin(ctx, request)
	}
	if strings.HasPrefix(request.ExportType, "pwntools") {
		return u.convertToPwntools(ctx, request)
	}
	var path string
	if strings.EqualFold(request.ExportType, "requests") {
		path = "convert_to_requests"
	} else if strings.EqualFold(request.ExportType, "curl") {
		path = "convert_to_curl"
	}

	res, err := req.Post[dto.ConvertToExploitResponse](
		fmt.Sprintf("http://altron.converter.loc:%d/%s", u.cfg.AltronConverterPort, path),
		dto.ConvertPacketToExploitRequest{
			Base64Str: request.Packet.Payload,
		},
	)
	if err != nil {
		return nil, err
	}

	return &spec.ConvertToExploitResponse{
		Exploit: res.Result,
	}, nil
}

func (u *ConversionUseCase) convertToAsciiBytes(ctx context.Context, request *spec.ConvertPacketToExploitRequest) (*spec.ConvertToExploitResponse, error) {
	decodedPayload, err := utils.FromBase64ToBytes(request.Packet.Payload)
	if err != nil {
		return nil, err
	}
	exploit := "payload = b''\n"

	transformedString := utils.TransformToAsciiBytes(decodedPayload)
	exploit += fmt.Sprintf("payload += b'%s'\n", strings.ReplaceAll(transformedString, "'", "\\'"))

	return &spec.ConvertToExploitResponse{
		Exploit: exploit,
	}, nil
}

func (u *ConversionUseCase) convertToPlugin(ctx context.Context, request *spec.ConvertPacketToExploitRequest) (*spec.ConvertToExploitResponse, error) {
	templateFile, err := template.ParseFiles("/core/templates/plugin.tmpl")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err := templateFile.Execute(&buf, request.Packet); err != nil {
		return nil, err
	}

	return &spec.ConvertToExploitResponse{
		Exploit: buf.String(),
	}, nil
}

func (u *ConversionUseCase) convertSessionToPwntools(ctx context.Context, request *spec.ConvertSessionToExploitRequest) (*spec.ConvertToExploitResponse, error) {
	session := request.Session
	packets := make([]spec.Packet, 0)
	for _, packet := range session.Packets {
		decodedPayload, err := utils.FromBase64ToBytes(packet.Payload)
		if err != nil {
			return nil, err
		}
		payload := strings.ReplaceAll(utils.TransformToAsciiBytes(decodedPayload), "'", "\\'")
		packets = append(packets, spec.Packet{
			Payload:   payload,
			IsRequest: packet.IsRequest,
		})
	}
	exploit := ""
	reqToVar := make(map[int]string, 0)
	for i := 0; i < len(packets); i++ {
		packet := packets[i]
		if packet.IsRequest {
			varName, ok := reqToVar[i]
			if ok {
				exploit += fmt.Sprintf("s.send(%s)\n", varName)
			} else {
				exploit += fmt.Sprintf("s.send(b'%s')\n", packet.Payload)
			}
		} else {
			var receiveCommand string
			if strings.EqualFold(request.ExportType, "pwntools_recvrepeat") {
				receiveCommand = "recvrepeat(timeout=1)"
			} else {
				startIndex := len(packet.Payload) - 10
				if startIndex < 0 {
					startIndex = 0
				}
				// check for hex cutting
				for idx := startIndex; idx >= 0; idx-- {
					if byte(packet.Payload[idx]) == '\\' {
						startIndex = idx
						break
					}
				}

				receiveCommand = fmt.Sprintf("recvuntil(b'%s')", packet.Payload[startIndex:])
			}
			exploit += fmt.Sprintf("s.%s\n", receiveCommand)
		}
	}
	templateFile, err := template.ParseFiles("/core/templates/pwntools.tmpl")
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := templateFile.Execute(&buf, ExploitTemplatePayload{
		AltronHost:  u.cfg.AltronHost,
		ServicePort: session.ServerPort,
		Exploit:     exploit,
	}); err != nil {
		return nil, err
	}

	return &spec.ConvertToExploitResponse{
		Exploit: buf.String(),
	}, nil
}

func (u *ConversionUseCase) convertToPwntools(ctx context.Context, request *spec.ConvertPacketToExploitRequest) (*spec.ConvertToExploitResponse, error) {
	decodedPayload, err := utils.FromBase64ToBytes(request.Packet.Payload)
	if err != nil {
		return nil, err
	}
	payload := strings.ReplaceAll(utils.TransformToAsciiBytes(decodedPayload), "'", "\\'")
	exploit := fmt.Sprintf("s.send(b'%s')\n", payload)
	templateFile, err := template.ParseFiles("/core/templates/pwntools.tmpl")
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := templateFile.Execute(&buf, ExploitTemplatePayload{
		ServicePort: request.ServicePort,
		Exploit:     exploit,
		AltronHost:  u.cfg.AltronHost,
	}); err != nil {
		return nil, err
	}

	return &spec.ConvertToExploitResponse{
		Exploit: buf.String(),
	}, nil
}

func (u *ConversionUseCase) ExtractFiles(ctx context.Context, request *spec.ExtractFilesFromPacketRequest) (*spec.ExtractFilesResponse, error) {
	zipFilename := fmt.Sprintf("/files/%s_%d.extracted.zip", request.SessionID, request.PacketNumber)
	defer os.Remove(zipFilename)

	if err := u.sftpClient.Download(zipFilename); err != nil {
		raw, err := base64.StdEncoding.DecodeString(request.Packet.Payload)
		if err != nil {
			return nil, err
		}

		u.log.Infof("Extracting files from packet %s_%d", request.SessionID, request.PacketNumber)

		binFilename := fmt.Sprintf("%s_%d", request.SessionID, request.PacketNumber)
		bin, err := os.Create(binFilename)
		if err != nil {
			return nil, err
		}
		if _, err := bin.Write(raw); err != nil {
			return nil, err
		}
		bin.Close()

		cmd := exec.Command("binwalk", "-D", ".", "-C", "/", binFilename, "--run-as=root")

		u.log.Infoln("binwalk", "-D", ".", "-C", "/", binFilename, "--run-as=root")

		if err := cmd.Run(); err != nil {
			return nil, err
		}

		u.log.Infof("Extracted files from packet %s_%d", request.SessionID, request.PacketNumber)

		if err := os.Remove(binFilename); err != nil {
			return nil, err
		}
		extractedDir := "_" + binFilename + ".extracted"

		if _, err := os.Stat(extractedDir); os.IsNotExist(err) {
			return nil, models.ErrorNoExtractedFiles
		}	

		u.log.Infof("Zipping extracted files from packet %s_%d", request.SessionID, request.PacketNumber)

		if err := u.zipDirectory(extractedDir, zipFilename); err != nil {
			return nil, err
		}

		u.log.Infof("Zipped extracted files from packet %s_%d", request.SessionID, request.PacketNumber)

		if err := os.RemoveAll(extractedDir); err != nil {
			return nil, err
		}

		u.log.Infof("Uploading zipped extracted files from packet %s_%d", request.SessionID, request.PacketNumber)

		if err := u.sftpClient.Upload(zipFilename); err != nil {
			return nil, err
		}
	}

	raw, err := os.ReadFile(zipFilename)
	if err != nil {
		return nil, err
	}

	b64Raw := base64.StdEncoding.EncodeToString(raw)
	return &spec.ExtractFilesResponse{
		Data: b64Raw,
	}, nil
}

func (u *ConversionUseCase) zipDirectory(dir, targetPath string) error {
	zipFile, err := os.Create(targetPath)
	if err != nil {
		u.log.Errorf("Error creating zip file %s: %v", targetPath, err)
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			u.log.Errorf("Error walking directory %s: %v", dir, err)
			return err
		}

		relPath, err := filepath.Rel(dir, filePath)
		if err != nil {
			u.log.Errorf("Error getting relative path %s: %v", filePath, err)
			return err
		}

		fileInZip, err := zipWriter.Create(relPath)
		if err != nil {
			u.log.Errorf("Error creating file in zip %s: %v", relPath, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			u.log.Errorf("Error opening file %s: %v", filePath, err)
			return err
		}
		defer file.Close()

		_, err = io.Copy(fileInZip, file)
		if err != nil {
			u.log.Errorf("Error copying file to zip %s: %v", filePath, err)
			return err
		}

		return nil
	})
}
