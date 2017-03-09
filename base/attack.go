package base

import "time"

// Attacker is our basic interface that allows us to be generic
type Attacker interface {
	Attack(...Validator) error
}

// BasicAttackerConfig is some fields that every attacker has
type BasicAttackerConfig struct {
	// NumThreads is the number of threads we use
	NumThreads uint
	// AttackDuration how long the attack should go for
	AttackDuration time.Duration
}

// StatefulAttackerConfig has connections and timeouts associated with them
type StatefulAttackerConfig struct {
	BasicAttackerConfig
	// NumConnections is the limit of connections we should have on each thread
	NumConnections uint
	// Timeout is our connection timeout
	Timeout time.Duration
	// ErrorThreshold is the percent of errors we are willing to encounter - -1 means unlimited
	ErrorThreshold int
}
