// +build linux

package networking

import "golang.org/x/sys/unix"

// SockHeaderIncl calls setsockopt IP_HDRINCL on our socket
func SockHeaderIncl(sock int) error {
	return unix.SetsockoptInt(sock, unix.IPPROTO_IP, unix.IP_HDRINCL, 1)
}
