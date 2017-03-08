package networking

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"golang.org/x/sys/unix"
)

func randomPort() uint16 {
	return binary.LittleEndian.Uint16(randomBytes(2))
}

func randomHighPort() uint16 {
	portNum := randomPort()
	if portNum < 1024 {
		portNum += (randomPort() % 1024) + 1024
	}
	return portNum
}

func randomBytes(numBytes int) []byte {
	b := make([]byte, numBytes)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return b
}

// htons converts a short (uint16) from host-to-network byte order.
// Thanks to mikioh for this neat trick:
// https://github.com/mikioh/-stdyng/blob/master/afpacket.go
func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

// BindToInterface binds a socket to an interface
func BindToInterface(fd int, ifaceName string) error {
	ifIndex := 0
	// An empty string here means to listen to all interfaces
	if ifaceName != "" {
		iface, err := net.InterfaceByName(ifaceName)
		if err != nil {
			return fmt.Errorf("InterfaceByName: %v", err)
		}
		ifIndex = iface.Index
	}
	s := &unix.SockaddrLinklayer{
		Protocol: htons(uint16(unix.ETH_P_ALL)),
		Ifindex:  ifIndex,
	}
	return unix.Bind(fd, s)
}

// CreateRawSocket returns a raw IP socket
func CreateRawSocket() (int, error) {
	return unix.Socket(unix.AF_INET, unix.SOCK_RAW, unix.IPPROTO_RAW)
}
