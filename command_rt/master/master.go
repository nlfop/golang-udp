package main

//go run master/master.go 1234 1235

import (
	"context"
	"fmt"
	"net"
	"os"
	"udp_connect/command_rt/handles/command"
	transmit "udp_connect/command_rt/handles/pkg"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}
	PORTData := ":" + arguments[1]
	PORTComm := ":" + arguments[2]
	fmt.Printf("%s - DataPort, %s - CommandsPort\n", PORTData, PORTComm)
	sData, err := net.ResolveUDPAddr("udp4", PORTData)

	if err != nil {
		fmt.Println(err)
		return
	}

	sComm, err := net.ResolveUDPAddr("udp4", PORTComm)

	if err != nil {
		fmt.Println(err)
		return
	}

	connectionData, err := net.ListenUDP("udp4", sData)
	if err != nil {
		fmt.Println(err)
		return
	}

	connectionComm, err := net.ListenUDP("udp4", sComm)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connectionData.Close()
	defer connectionComm.Close()

	buffer := make([]byte, 1024)
	ctx, cancel := context.WithCancel(context.Background())
	sAddrData, err := net.ResolveUDPAddr("udp4", "127.0.0.1:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		n, _, _ := connectionComm.ReadFromUDP(buffer)

		comm, err := command.CommandTrim(buffer[0:n])
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch comm {
		case "STOP_FLOW":
			cancel()
			fmt.Println("Exiting UDP server!")
			ctx, cancel = context.WithCancel(context.Background())
		case "START_FLOW":
			go transmit.TransmitStructure(ctx, cancel, connectionData, sAddrData)
		case "START_ONCE":
			go transmit.TransmitStructureOnce(connectionData, sAddrData)

		}

	}
}
