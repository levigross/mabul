package http

import (
	"fmt"
	"sync"
)

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
