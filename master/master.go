package main

//go run master/master.go 1234

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	transmit "udp_connect/handles/pkg"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}
	PORT := ":" + arguments[1]

	s, err := net.ResolveUDPAddr("udp4", PORT)

	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()
	buffer := make([]byte, 1024)
	ctx, cancel := context.WithCancel(context.Background())
	for {
		n, addr, _ := connection.ReadFromUDP(buffer)
		fmt.Print("-> ", string(buffer[0:n-1]))

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			cancel()
			fmt.Println("Exiting UDP server!")
			ctx, cancel = context.WithCancel(context.Background())
			// return
		}
		if strings.TrimSpace(string(buffer[0:n])) == "START" {

			go transmit.TransmitStructure(ctx, connection, addr)

		}

	}
}
