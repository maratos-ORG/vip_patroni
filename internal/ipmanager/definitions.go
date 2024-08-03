package ipmanager

import (
	"net"
	"net/netip"
	"sync"

	arp "github.com/mdlayher/arp"
)

type IPManager struct {
	Configurer   ipConfigurer
	states       <-chan bool
	CurrentState bool
	stateLock    sync.Mutex
	recheck      *sync.Cond
}

// IPConfiguration holds the configuration for vip_patroni
type IPConfiguration struct {
	VIP        netip.Addr
	Netmask    net.IPMask
	Iface      net.Interface
	RetryNum   int
	RetryAfter int
}

type ipConfigurer interface {
	QueryAddress() bool
	configureAddress() bool
	deconfigureAddress() bool
	GetCIDR() string
	cleanupArp()
}

// BasicConfigurer can be used to enable vip-management on nodes
// that handle their own network connection, in setups where it is
// sufficient to add the virtual ip using `ip addr add ...` .
// After adding the virtual ip to the specified interface,
// a gratuitous ARP package is sent out to update the tables of
// nearby routers and other devices.
type BasicConfigurer struct {
	*IPConfiguration
	arpClient *arp.Client
}
