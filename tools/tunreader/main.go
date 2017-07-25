package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

var (
	flagTun = flag.String("tun", "tuntest", "Tunnel device name")
)

func main() {
	iface, err := water.New(water.Config{
		DeviceType:             water.TUN,
		PlatformSpecificParams: water.PlatformSpecificParams{Name: *flagTun}})
	if err != nil {
		fmt.Println("Error water.New:", err)
		os.Exit(1)
	}

	link, err := netlink.LinkByName(*flagTun)
	if err != nil {
		fmt.Println("Error netlink.LinkByName:", err)
		os.Exit(1)
	}

	err = netlink.LinkSetUp(link)
	if err != nil {
		fmt.Println("Error nelink.LinkSetUp:", err)
		os.Exit(1)
	}

	buffer := make([]byte, 1<<16)
	for {
		_, err = iface.Read(buffer)
		if err != nil {
			fmt.Println("Error Read:", err)
		}
	}
}
