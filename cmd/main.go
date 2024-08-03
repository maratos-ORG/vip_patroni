package main

import (
	"context"
	"net"
	"net/netip"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"vip_patroni/internal/api"
	"vip_patroni/internal/checker"
	"vip_patroni/internal/config"
	"vip_patroni/internal/ipmanager"
	log "vip_patroni/internal/logging"
)

var (
	build = "0.0.3"
)

func main() {
	conf, err := config.NewConfig(build)
	if err != nil {
		log.Fatal("unable to generate config (pkg->config)")
	}
	log.Level(conf.LogLevel)
	vip := netip.MustParseAddr(conf.IP)
	vipMask := GetMask(vip, conf.Mask)
	netIface := GetNetIface(conf.Iface)
	states := make(chan bool)
	mainCtx, cancel := context.WithCancel(context.Background())

	manager, err := ipmanager.NewIPManager(
		&ipmanager.IPConfiguration{
			VIP:        vip,
			Netmask:    vipMask,
			Iface:      *netIface,
			RetryNum:   conf.RetryNum,
			RetryAfter: conf.RetryAfter,
		},
		states,
	)

	if err != nil {
		log.Fatal("problems with generating the virtual ip manage")
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Warning("Received exit signal")
		cancel()
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := checker.InitChecker(mainCtx, conf, states)
		if err != nil && err != context.Canceled {
			log.Fatal("Role checker returned the following error")
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		manager.SyncStates(mainCtx, states)
		wg.Done()
	}()

	go func() {
		api.NewAPI(conf, manager)
	}()

	wg.Wait()
}

// Validate value of interface name
func GetNetIface(iface string) *net.Interface {
	netIface, err := net.InterfaceByName(iface)
	if err != nil {
		log.Fatal("obtaining the interface raised an error")
	}
	return netIface
}

// Validate value of netmask1
func GetMask(vip netip.Addr, mask int) net.IPMask {
	if vip.Is4() { //IPv4
		if mask > 0 || mask < 33 {
			return net.CIDRMask(mask, 32)
		}
		var ip net.IP = vip.AsSlice()
		return ip.DefaultMask()
	}
	return net.CIDRMask(mask, 128) //IPv6
}
