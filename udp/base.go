package udp

import (
	"time"

	"github.com/levigross/mabul/base"
	"go.uber.org/zap"
)

// AttackType ..
type AttackType int

const (
	// SNMP ...
	SNMP AttackType = iota
	// DNS ...
	DNS
	// NTP ...
	NTP
	// SSDP ...
	SSDP
)

// DefaultPort will give you the default port of the protocol you wish to attack
func (a AttackType) DefaultPort() uint16 {
	switch a {
	case SNMP:
		return 161
	case NTP:
		return 123
	case SSDP:
		return 1900
	case DNS:
		return 53
	default: // This should never happen
		return 0
	}
}

// AttackConfig ...
type AttackConfig struct {
	base.BasicAttackerConfig
	// AttackType the type of attack we want
	AttackType AttackType
}

// Attacker is our basic attacker that send UDP
type Attacker struct {
	Config *AttackConfig
	Target base.Target
	Log    *zap.SugaredLogger

	UDPAttacker DerAttacker
}

var _ base.Attacker = &Attacker{}

// Attack will send out UDP attacks
func (a *Attacker) Attack(v ...base.Validator) error {
	if err := base.Validate(v...); err != nil {
		return err
	}

	switch a.Config.AttackType {
	case SNMP:
		a.UDPAttacker = &SNMPFlood{}
		if a.Target.DstPort == 0 {
			a.Target.DstPort = SNMP.DefaultPort()
		}
	}
	a.UDPAttacker.SetAttacker(a)

	select {
	case <-time.After(a.Config.AttackDuration):
		return nil
	case err := <-a.attack():
		return err
	}
}

func (a *Attacker) attack() <-chan error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- a.UDPAttacker.Attack()
	}()
	return errChan
}
