//go:build linux
// +build linux

package ipmanager

import (
	"errors"
	"net"
	"os/exec"
	"strings"
	"time"
	log "vip_patroni/internal/logging"

	arp "github.com/mdlayher/arp"
)

const (
	arpRequestOp = 1
	arpReplyOp   = 2
)

var (
	ethernetBroadcast = net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

// configureAddress assigns virtual IP address
func (c *BasicConfigurer) configureAddress() bool {
	if c.arpClient == nil {
		err := c.createArpClient()
		if err != nil {
			// log.Fatalf("Couldn't create an Arp client: %s", err)
			log.Fatal("Couldn't create an Arp client: %s", err)
		}
	}

	// log.Printf("Configuring address %s on %s", c.getCIDR(), c.Iface.Name)
	log.Info("Configuring address %s on %s", c.GetCIDR(), c.Iface.Name)
	result := c.runAddressConfiguration("add")

	if result {
		// For now it is save to say that also working even if a
		// gratuitous arp message could not be send but logging an
		// errror should be enough.
		_ = c.arpSendGratuitous()
	}

	return result
}

// deconfigureAddress drops virtual IP address
func (c *BasicConfigurer) deconfigureAddress() bool {
	// log.Printf("Removing address %s on %s", c.getCIDR(), c.Iface.Name)
	log.Info("Removing address %s on %s", c.GetCIDR(), c.Iface.Name)
	return c.runAddressConfiguration("delete")
}

func (c *BasicConfigurer) runAddressConfiguration(action string) bool {
	cmd := exec.Command("ip", "addr", action,
		c.GetCIDR(),
		"dev", c.Iface.Name)
	output, err := cmd.CombinedOutput()

	switch err.(type) {
	case *exec.ExitError:
		log.Info("Got error %s", output)
		return false
	}
	if err != nil {
		log.Info("Error running ip address %s %s on %s: %s", action, c.VIP, c.Iface.Name, err)
		return false
	}
	return true
}

func (c *BasicConfigurer) createArpClient() error {
	var err error
	var arpClient *arp.Client
	for i := 0; i < c.RetryNum; i++ {
		arpClient, err = arp.Dial(&c.Iface)
		if err != nil {
			log.Info("Problems with producing the arp client: %s", err)
		} else {
			break
		}
		time.Sleep(time.Duration(c.RetryAfter) * time.Millisecond)
	}
	if err != nil {
		log.Info("too many retries")
		return err
	}
	c.arpClient = arpClient
	return nil
}

// sends a gratuitous ARP request and reply
func (c *BasicConfigurer) arpSendGratuitous() error {
	/* While RFC 2002 does not say whether a gratuitous ARP request or reply is preferred
	 * to update ones neighbours' MAC tables, the Wireshark Wiki recommends sending both.
	 *		https://wiki.wireshark.org/Gratuitous_ARP
	 * This site also recommends sending a reply, as requests might be ignored by some hardware:
	 *		https://support.citrix.com/article/CTX112701
	 */
	gratuitousReplyPackage, err := arp.NewPacket(
		arpReplyOp,
		c.Iface.HardwareAddr,
		c.VIP.AsSlice(),
		c.Iface.HardwareAddr,
		c.VIP.AsSlice(),
	)
	if err != nil {
		log.Info("Gratuitous arp reply package is malformed: %s", err)
		return err
	}

	/* RFC 2002 specifies (in section 4.6) that a gratuitous ARP request
	 * should "not set" the target Hardware Address (THA).
	 * Since the arp package offers no option to leave the THA out, we specify the Zero-MAC.
	 * If parsing that fails for some reason, we'll just use the local interface's address.
	 * The field is probably ignored by the receivers' implementation anyway.
	 */
	arpRequestDestMac, err := net.ParseMAC("00:00:00:00:00:00")
	if err != nil {
		// not entirely RFC-2002 conform but better then nothing.
		arpRequestDestMac = c.Iface.HardwareAddr
	}

	gratuitousRequestPackage, err := arp.NewPacket(
		arpRequestOp,
		c.Iface.HardwareAddr,
		c.VIP.AsSlice(),
		arpRequestDestMac,
		c.VIP.AsSlice(),
	)
	if err != nil {
		log.Info("Gratuitous arp request package is malformed: %s", err)
		return err
	}

	for i := 0; i < c.RetryNum; i++ {
		errReply := c.arpClient.WriteTo(gratuitousReplyPackage, ethernetBroadcast)
		if err != nil {
			log.Info("Couldn't write to the arpClient: %s", errReply)
		} else {
			log.Info("Sent gratuitous ARP reply")
		}

		errRequest := c.arpClient.WriteTo(gratuitousRequestPackage, ethernetBroadcast)
		if err != nil {
			log.Info("Couldn't write to the arpClient: %s", errRequest)
		} else {
			log.Info("Sent gratuitous ARP request")
		}

		if errReply != nil || errRequest != nil {
			/* If something went wrong while sending the packages, we'll recreate the ARP client for the next try,
			 * to avoid having a stale client that gives "network is down" error.
			 */
			err = c.createArpClient()
		} else {
			//TODO: think about whether to leave this out to achieve simple repeat sending of GARP packages
			break
		}
		time.Sleep(time.Duration(c.RetryAfter) * time.Millisecond)
	}
	if err != nil {
		log.Info("too many retries")
		return err
	}

	return nil
}

func newBasicConfigurer(config *IPConfiguration) (*BasicConfigurer, error) {
	c := &BasicConfigurer{IPConfiguration: config}
	if c.Iface.HardwareAddr == nil || c.Iface.HardwareAddr.String() == "00:00:00:00:00:00" {
		return nil, errors.New(`Cannot run vip-manager on the loopback device
as its hardware address is the local address (00:00:00:00:00:00),
which prohibits sending of gratuitous ARP messages`)
	}
	return c, nil
}

// queryAddress returns if the address is assigned
func (c *BasicConfigurer) QueryAddress() bool {
	iface, err := net.InterfaceByName(c.Iface.Name)
	if err != nil {
		return false
	}
	addresses, err := iface.Addrs()
	if err != nil {
		return false
	}
	for _, address := range addresses {
		if strings.Contains(address.String(), c.VIP.String()) {
			return true
		}
	}
	return false
}

func (c *BasicConfigurer) cleanupArp() {
	if c.arpClient != nil {
		c.arpClient.Close()
	}
}
