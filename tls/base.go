package tls

import (
	"github.com/levigross/gboringssltls"
	"github.com/levigross/mabul/base"
)

// AttackType the type of TLS attack you want to set
type AttackType int

const (
	// ClientHelloFlood will flood a connection with Hello requests
	ClientHelloFlood AttackType = iota
	// MassConnect will connect a bunch of times and do nothing
	MassConnect
	// Politician will do anything to connect but once connected, will
	// do what it wants
	Politician
)

// AttackConfig helps us configure our TLS attacker
type AttackConfig struct {
	base.StatefulAttackerConfig

	// AttackType is our attack type
	AttackType AttackType

	// RandomServerName will present a random server name when sending a ClientHello
	RandomServerName bool

	// ExpensiveCiphers will ensure we use only the most expensive ciphers
	ExpensiveCiphers bool

	// OnlyRSA will ensure we only use expensive RSA ciphers
	OnlyRSA bool

	// TLSVersion which TLS version do you want to use?
	TLSVersion string

	// PreferServerCiphers allow the server to dictate the ciphers
	PreferServerCiphers bool

	// DisableVerification will turn off TLS cert verification
	DisableVerification bool

	// ServerName is an option you set when you don't want the domain within base.Target
	// to be the domain you wish to target. Use of this and RandomServerName is mutually exclusive
	ServerName string

	// ProtocolBugs are some bugs that we can just inject into the protocol
	ProtocolBugs gboringssltls.ProtocolBugs
}
