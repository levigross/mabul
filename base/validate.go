package base

import (
	"errors"
	"fmt"
	"net"
)

// Validator allows us to remove types from the validator
type Validator interface {
	Validate() error
}

var _ Validator = (*Target)(nil)

// Target has some basic details about our target
type Target struct {
	DomainName    string
	IPAddress     net.IP
	DstPort       uint16 // 0 means randomize
	InterfaceName string
}

// Validate will validate the target information
func (t *Target) Validate() error {
	if t.DomainName == "" && t.IPAddress == nil {
		return errors.New("You must choose either a domain or IP to attack")
	}
	if t.DomainName != "" {
		ips, err := net.LookupIP(t.DomainName)
		if err != nil {
			return err
		}
		t.IPAddress = ips[0]
	}

	if t.InterfaceName != "" {
		if _, err := net.InterfaceByName(t.InterfaceName); err != nil {
			return fmt.Errorf("InterfaceByName: %v", err)
		}
	}
	return nil
}

// Validate will validate many validators
func Validate(v ...Validator) error {
	for _, f := range v {
		if err := f.Validate(); err != nil {
			return err
		}
	}
	return nil
}
