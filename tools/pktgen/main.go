package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	ipFlag   = flag.String("ip", "", "destination IP")
	portFlag = flag.Uint("port", 0, "destination Port")
	sizeFlag = flag.Uint("size", 64, "size of UDP payload, in bytes")
)

var (
	ip   net.IP
	port uint
	size uint
)

func gen() {
	address := &net.UDPAddr{IP: ip, Port: int(port)}
	conn, err := net.DialUDP("udp4", nil, address)
	if err != nil {
		fmt.Printf("Dial error (%v)\n", err)
		return
	}

	buffer := make([]byte, size)
	for {
		conn.Write(buffer)
		// Don't even care about errors here
	}
}

func main() {
	flag.Parse()

	ip = net.ParseIP(*ipFlag)
	if ip == nil {
		fmt.Printf("Unable to parse IP %s\n", *ipFlag)
		os.Exit(1)
	}

	port = *portFlag
	if port == 0 || port >= 1<<16 {
		fmt.Printf("Invalid port number %d\n", *portFlag)
		os.Exit(1)
	}

	size = *sizeFlag

	gen()
}
