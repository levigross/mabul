package base

// Attacker is our basic interface that allows us to be generic
type Attacker interface {
	Attack(...Validator) error
}
