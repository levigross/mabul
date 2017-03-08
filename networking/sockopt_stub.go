// +build !linux

package networking

import "errors"

// SockHeaderIncl ...
func SockHeaderIncl(sock int) error {
	return errors.New("This option is not supported on your OS")
}
