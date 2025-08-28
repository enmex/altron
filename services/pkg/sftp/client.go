package sftp

import (
	"io"
	"os"
	"sync"

	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	cfg    *Config
	client *sftp.Client
	mut    *sync.Mutex
}

func NewClient(cfg *Config) (*Client, error) {
	sshClient, err := ssh.Dial("tcp", cfg.Server, &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}

	client, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	return &Client{
		cfg:    cfg,
		client: client,
		mut:    &sync.Mutex{},
	}, nil
}

func (c *Client) Download(filePath string) error {
	c.mut.Lock()
	defer c.mut.Unlock()
	localFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer localFile.Close()
	remoteFile, err := c.client.Open(filePath[1:])
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	_, err = io.Copy(localFile, remoteFile)
	return err
}

func (c *Client) List(path string) ([]string, error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	filesInfo, err := c.client.ReadDir(path)
	if err != nil {
		return nil, err
	}
	filenames := make([]string, 0, len(filesInfo))
	for _, fileInfo := range filesInfo {
		if fileInfo.IsDir() {
			continue
		}
		filenames = append(filenames, fileInfo.Name())
	}
	return filenames, nil
}

func (c *Client) Delete(filepath string) error {
	c.mut.Lock()
	defer c.mut.Unlock()
	return c.client.Remove(filepath[1:])
}

func (c *Client) MakeDir(path string) error {
	c.mut.Lock()
	defer c.mut.Unlock()
	return c.client.Mkdir(path)
}

func (c *Client) Upload(filePath string) error {
	c.mut.Lock()
	defer c.mut.Unlock()
	remoteFile, err := c.client.Create(filePath[1:])
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	localFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	_, err = io.Copy(remoteFile, localFile)
	return err
}
