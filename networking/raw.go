package networking

import (
	"net"

	"golang.org/x/sys/unix"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// EndPoint is what we use to send raw packets on our endpoint
type EndPoint struct {
	Src              net.IP
	DstPort          uint16
	SrcPort          uint16
	RandomizeSrcPort bool
	RandomizeDstPort bool

	target int
}

// NewEndPoint returns an attack endpoint
func NewEndPoint(src net.IP, dstPort uint16,
	networkInterface string) (*EndPoint, error) {
	sock, err := CreateRawSocket()
	if err != nil || sock < 0 {
		return nil, err
	}

	if err := SockHeaderIncl(sock); err != nil {
		return nil, err
	}

	return &EndPoint{Src: src, DstPort: dstPort, target: sock}, nil
}

// SendUDPPacket will send a UDP packet
func (e *EndPoint) SendUDPPacket(payload []byte, Dst net.IP) error {
	ip := &layers.IPv4{
		TTL:      255,
		Version:  4,
		SrcIP:    e.Src,
		DstIP:    Dst,
		Protocol: layers.IPProtocolUDP,
	}
	udp := &layers.UDP{
		SrcPort: layers.UDPPort(randomHighPort()),
		DstPort: layers.UDPPort(e.DstPort),
	}
	buf := gopacket.NewSerializeBuffer() // TODO: Reuse buffers
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	if err := udp.SetNetworkLayerForChecksum(ip); err != nil {
		return err
	}
	if err := gopacket.SerializeLayers(buf, opts, ip, udp, gopacket.Payload(payload)); err != nil {
		return err
	}
	return e.SendPacket(buf.Bytes())
}

// SendPacket send the packet in a unix agnostic way
func (e *EndPoint) SendPacket(b []byte) error {
	return unix.Sendto(e.target, b, 0, &unix.SockaddrInet4{})
}
