package main

//go run slave/slave.go 127.0.0.1:1234 127.0.0.1:1235 127.0.0.1:8000 127.0.0.1:8001

// 68 00 10 10 00 00 00 - старт периодической передачи
// 68 00 20 01 00 00 00 - стоп периодической передачи

// 68000701000000 кс старт
// 68000702000000 кс стоп

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"
	"udp_connect/command_rt/handles/command"
	transmit "udp_connect/command_rt/handles/pkg"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a host:port string")
		return
	}
	CONNECTData := arguments[1]
	CONNECTComm := arguments[2]
	addrData := arguments[3]
	addrComm := arguments[4]

	sData, err := net.ResolveUDPAddr("udp4", CONNECTData)
	if err != nil {
		fmt.Println(err)
		return
	}
	sComm, err := net.ResolveUDPAddr("udp4", CONNECTComm)
	if err != nil {
		fmt.Println(err)
		return
	}

	sAddrData, err := net.ResolveUDPAddr("udp4", addrData)
	if err != nil {
		fmt.Println(err)
		return
	}
	sAddrComm, err := net.ResolveUDPAddr("udp4", addrComm)
	if err != nil {
		fmt.Println(err)
		return
	}

	cData, err := net.DialUDP("udp4", sAddrData, sData)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cData.Close()

	cComm, err := net.DialUDP("udp4", sAddrComm, sComm)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s - Data, %s - Commands\n", cData.RemoteAddr().String(), cComm.RemoteAddr().String())
	fmt.Println("Enter 7 bytes without spaces and checksum (it counts automatically)")
	defer cData.Close()
	defer cComm.Close()
	var ctx context.Context
	var cancel context.CancelFunc

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		data, _ := reader.ReadString('\n')
		n := len(data)

		str, _ := hex.DecodeString(data[0 : n-1])
		str = transmit.CountCheckSum(str)
		_, err = cComm.Write(str)
		if err != nil {
			fmt.Println("--error")
			return
		}

		comm, err := command.CommandTrim(str)
		if err != nil {
			fmt.Println(err)
			continue
		}
		switch comm {
		case "STOP_FLOW":

			time.Sleep(2000 * time.Millisecond)
			fmt.Println("Exiting UDP client flow data!")

		case "START_FLOW":
			ctx, cancel = context.WithCancel(context.Background())
			go transmit.ReceiveStructure(cData, ctx, cancel)
			// case "START_ONCE":
			// 	go transmit.ReceiveStructureOnce(cData)
		}

	}
}
