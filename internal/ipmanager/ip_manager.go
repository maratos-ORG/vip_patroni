package ipmanager

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
	log "vip_patroni/internal/logging"
)

// NewIPManager returns a new instance of IPManager
func NewIPManager(config *IPConfiguration, states <-chan bool) (m *IPManager, err error) {
	m = &IPManager{
		states:       states,
		CurrentState: false,
	}
	m.recheck = sync.NewCond(&m.stateLock)

	m.Configurer, err = newBasicConfigurer(config)
	if err != nil {
		m = nil
	}
	return
}

// SyncStates implements states synchronization
func (m *IPManager) SyncStates(ctx context.Context, states <-chan bool) {
	ticker := time.NewTicker(60 * time.Second)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		m.applyLoop(ctx)
		wg.Done()
	}()

	for {
		select {
		case newState := <-states:
			m.stateLock.Lock()
			if m.CurrentState != newState {
				m.CurrentState = newState
				m.recheck.Broadcast()
			}
			m.stateLock.Unlock()
		case <-ticker.C:
			m.recheck.Broadcast()
		case <-ctx.Done():
			m.recheck.Broadcast()
			wg.Wait()
			m.Configurer.cleanupArp()
			return
		}
	}
}

func (m *IPManager) applyLoop(ctx context.Context) {
	timeout := 0
	for {
		select {
		case <-ctx.Done():
			m.Configurer.deconfigureAddress()
			return
		case <-time.After(time.Duration(timeout) * time.Second):
			actualState := m.Configurer.QueryAddress()
			m.stateLock.Lock()
			desiredState := m.CurrentState
			log.Info("IP address %s state is %t, desired %t", m.Configurer.GetCIDR(), actualState, desiredState)
			if actualState != desiredState {
				m.stateLock.Unlock()
				var configureState bool
				if desiredState {
					configureState = m.Configurer.configureAddress()
				} else {
					configureState = m.Configurer.deconfigureAddress()
				}
				if !configureState {
					log.Info("Error while acquiring virtual ip for this machine")
					//Sleep a little bit to avoid busy waiting due to the for loop.
					timeout = 10
				} else {
					timeout = 0
				}
			} else {
				// Wait for notification
				m.recheck.Wait()
				// Want to query actual state anyway, so unlock
				m.stateLock.Unlock()
			}
		}
	}
}

// getCIDR returns the CIDR composed from the given address and mask
func (c *IPConfiguration) GetCIDR() string {
	return fmt.Sprintf("%s/%d", c.VIP.String(), netmaskSize(c.Netmask))
}

func netmaskSize(mask net.IPMask) int {
	ones, bits := mask.Size()
	if bits == 0 {
		panic("Invalid mask")
	}
	return ones
}
