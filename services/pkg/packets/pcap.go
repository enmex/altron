package packets

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type PcapAnalyzer struct{}

func NewPcapAnalyzer() *PcapAnalyzer {
	return &PcapAnalyzer{}
}

func (pa *PcapAnalyzer) PcapFileHandler(filename string) (*pcap.Handle, error) {
	handler, err := pcap.OpenOffline(filename)
	if err != nil {
		return nil, err
	}
	return handler, nil
}

func (pa *PcapAnalyzer) PacketSource(handler *pcap.Handle) <-chan gopacket.Packet {
	return gopacket.NewPacketSource(handler, handler.LinkType()).Packets()
}
