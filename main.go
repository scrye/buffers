package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

func server() {
	listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: "/run/shm/test.sock", Net: "unix"})
	if err != nil {
		fmt.Println("ListenUnix", err)
		return
	}

	conn, err := listener.AcceptUnix()
	if err != nil {
		fmt.Println("AcceptUnix", err)
		return
	}

	buffer := make([]byte, 1500)
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Read", err)
			return
		}
	}
}

func client() {
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: "/run/shm/test.sock", Net: "unix"})
	if err != nil {
		fmt.Println("DialUnix", err)
		return
	}

	buffer := make([]byte, 100)
	for i := 0; i < 20; i++ {
		before := time.Now().UnixNano()
		//n, err := conn.Write(buffer)
		_, err := io.Copy(conn, bytes.NewReader(buffer))
		after := time.Now().UnixNano()
		if err != nil {
			fmt.Println("Write", err)
			return
		}
		fmt.Println("Time (us): ", (after-before)/1000)
	}
}

func main() {
	go server()
	time.Sleep(time.Second)
	client()
}
