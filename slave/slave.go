package main

//go run slave/slave.go 127.0.0.1:1234

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"
	"udp_connect/handles/command"
	transmit "udp_connect/handles/pkg"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a host:port string")
		return
	}
	CONNECT := arguments[1]

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	fmt.Println("Enter 7 bytes without spaces and checksum (it counts automatically)")
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		data, _ := reader.ReadString('\n')
		n := len(data)

		str, _ := hex.DecodeString(data[0 : n-1])
		str = transmit.CountCheckSum(str)
		_, err = c.Write(str)
		if err != nil {
			fmt.Println("--error")
			return
		}
		fmt.Println(str)
		comm, err := command.CommandTrim(str)
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch comm {
		case "STOP_FLOW":
			cancel()
			time.Sleep(2000 * time.Millisecond)
			fmt.Println("Exiting UDP client!")
			return
		case "START_FLOW":
			go transmit.ReceiveStructure(c, ctx)
		case "START_ONCE":
			go transmit.ReceiveStructureOnce(c)
		}

	}
}
