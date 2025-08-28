package models

import (
	"altron/common/models"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/xlab/treeprint"
)

type SessionTree struct {
	tcpSessionTimeout time.Duration
	size              atomic.Uint64
	servers           sync.Map //*utils.ConcurrentMap[uint16, *Server]
}

type Server struct {
	size       atomic.Uint64
	interfaces sync.Map //*utils.ConcurrentMap[string, *Interface]
}

type Interface struct {
	size    atomic.Uint64
	clients sync.Map //*utils.ConcurrentMap[string, *Clients]
}

type Client struct {
	FileName *string
	Session  *models.Session
}

func NewSessionTree(tcpSessionTimeout time.Duration) *SessionTree {
	return &SessionTree{
		tcpSessionTimeout: tcpSessionTimeout,
		servers:           sync.Map{},
		size:              atomic.Uint64{},
	}
}

func (t *SessionTree) AddServer(serverPort uint16) {
	if _, ok := t.servers.Load(serverPort); ok {
		return
	}
	t.size.Add(1)
	t.servers.LoadOrStore(serverPort, &Server{
		interfaces: sync.Map{},
		size:       atomic.Uint64{},
	})
}

func (t *SessionTree) AddInterface(serverPort uint16, iface string) error {
	v, ok := t.servers.Load(serverPort)
	if !ok {
		return ErrorServerNotFound
	}
	server := v.(*Server)
	server.size.Add(1)
	server.interfaces.LoadOrStore(iface, &Interface{
		clients: sync.Map{},
		size:    atomic.Uint64{},
	})
	return nil
}

func (t *SessionTree) AddClientHost(
	serverPort uint16,
	iface string,
	clientHost string,
	protocol string,
	sentAt time.Time,
	TTL uint8,
	fileName *string,
	userAgent *string,
) error {
	interfaceNode, err := t.interfaceNode(serverPort, iface)
	if err != nil {
		return err
	}
	interfaceNode.size.Add(1)
	interfaceNode.clients.LoadOrStore(clientHost, &Client{
		FileName: fileName,
		Session: &models.Session{
			ID:              uuid.New(),
			Iface:           iface,
			ServerPort:      serverPort,
			ClientHost:      clientHost,
			Protocol:        protocol,
			SentAt:          sentAt,
			TTL:             TTL,
			Packets:         make([]*models.Packet, 0),
			ClientUserAgent: userAgent,
		},
	})
	return nil
}

func (t *SessionTree) Servers() []uint16 {
	keys := make([]uint16, 0)
	t.servers.Range(func(key, value any) bool {
		port := key.(uint16)
		keys = append(keys, port)
		return true
	})
	return keys
}

func (t *SessionTree) Interfaces(serverPort uint16) ([]string, error) {
	v, ok := t.servers.Load(serverPort)
	if !ok {
		return nil, ErrorInterfaceNotFound
	}
	interfaceNames := make([]string, 0)

	serverNode := v.(*Server)
	serverNode.interfaces.Range(func(key, value any) bool {
		interfaceNames = append(interfaceNames, key.(string))
		return true
	})

	return interfaceNames, nil
}

func (t *SessionTree) ClientHosts(serverPort uint16, iface string) ([]string, error) {
	interfaceNode, err := t.interfaceNode(serverPort, iface)
	if err != nil {
		return nil, err
	}
	clients := make([]string, 0)
	interfaceNode.clients.Range(func(key, value any) bool {
		clients = append(clients, key.(string))
		return true
	})
	return clients, nil
}

func (t *SessionTree) ClientExists(serverPort uint16, iface string, clientHost string) (bool, error) {
	_, err := t.clientHostNode(serverPort, iface, clientHost)
	if err != nil {
		if errors.Is(err, ErrorClientNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (t *SessionTree) ServerExists(serverPort uint16) bool {
	_, ok := t.servers.Load(serverPort)
	return ok
}

func (t *SessionTree) PortInterfaceExists(serverPort uint16, iface string) bool {
	_, err := t.interfaceNode(serverPort, iface)
	return err == nil
}

func (t *SessionTree) AddPacket(serverPort uint16, iface string, clientHost string, packet *models.Packet) error {
	clientHostNode, err := t.clientHostNode(serverPort, iface, clientHost)
	if err != nil {
		return err
	}
	if len(clientHostNode.Session.Packets) >= 1 && packet.SentAt.Sub(clientHostNode.Session.Packets[0].SentAt) > t.tcpSessionTimeout {
		return nil
	}

	clientHostNode.Session.Packets = append(clientHostNode.Session.Packets, packet)
	sort.Slice(clientHostNode.Session.Packets, func(i, j int) bool {
		return clientHostNode.Session.Packets[i].SentAt.Before(clientHostNode.Session.Packets[j].SentAt)
	})
	return nil
}

func (t *SessionTree) ClientSession(serverPort uint16, iface string, clientHost string) (*models.Session, *string, error) {
	clientHostNode, err := t.clientHostNode(serverPort, iface, clientHost)
	if err != nil {
		return nil, nil, err
	}
	return clientHostNode.Session, clientHostNode.FileName, nil
}

func (t *SessionTree) DeleteServer(serverPort uint16) {
	current := t.size.Load()
	if current != 0 {
		t.size.Swap(current - 1)
	}
	t.servers.Delete(serverPort)
}

func (t *SessionTree) DeleteClient(iface string, serverPort uint16, clientHost string) error {
	ifaceNode, err := t.interfaceNode(serverPort, iface)
	if err != nil {
		return err
	}
	current := ifaceNode.size.Load()
	if current != 0 {
		ifaceNode.size.Swap(current - 1)
	}
	ifaceNode.clients.Delete(clientHost)
	return nil
}

func (t *SessionTree) clientHostNode(serverPort uint16, iface string, clientHost string) (*Client, error) {
	interfaceNode, err := t.interfaceNode(serverPort, iface)
	if err != nil {
		return nil, err
	}
	v, ok := interfaceNode.clients.Load(clientHost)
	if !ok {
		return nil, ErrorClientNotFound
	}
	clientNode := v.(*Client)
	return clientNode, nil
}

func (t *SessionTree) interfaceNode(serverPort uint16, iface string) (*Interface, error) {
	v, ok := t.servers.Load(serverPort)
	if !ok {
		return nil, ErrorServerNotFound
	}
	serverNode := v.(*Server)
	v, ok = serverNode.interfaces.Load(iface)
	if !ok {
		return nil, ErrorInterfaceNotFound
	}

	interfaceNode := v.(*Interface)
	return interfaceNode, nil
}

func (t *SessionTree) String() string {
	tree := treeprint.New()

	for _, serverPort := range t.Servers() {
		v, ok := t.servers.Load(serverPort)
		if !ok {
			continue
		}
		serverNode := v.(*Server)
		if serverNode.size.Load() == 0 {
			tree.AddNode(fmt.Sprint(serverPort))
			continue
		}
		portBranch := tree.AddBranch(fmt.Sprint(serverPort))

		ifaces, err := t.Interfaces(serverPort)
		if err != nil {
			continue
		}
		for _, iface := range ifaces {
			interfaceNode, err := t.interfaceNode(serverPort, iface)
			if err != nil {
				continue
			}
			if interfaceNode.size.Load() == 0 {
				portBranch.AddNode(iface)
				continue
			}
			clientBranch := portBranch.AddBranch(iface)

			clients, err := t.ClientHosts(serverPort, iface)
			if err != nil {
				continue
			}
			for _, client := range clients {
				clientNode, err := t.clientHostNode(serverPort, iface, client)
				if err != nil {
					continue
				}
				session := clientNode.Session
				clientBranch.AddNode(
					fmt.Sprintf(
						"[%s]: %d packets, sent at %s",
						strings.Split(session.ClientHost, ":")[1],
						len(session.Packets),
						session.SentAt.Format("15:04:05"),
					),
				)
			}
		}
	}

	return fmt.Sprintf("Session tree:\n%s", tree.String())
}
