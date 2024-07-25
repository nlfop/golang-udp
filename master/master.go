package main

//go run master/master.go 1234

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

type Packet struct {
	PortTo   int
	Message  string
	NumFloat float32
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

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
	rand.Seed(time.Now().Unix())

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		fmt.Print("-> ", string(buffer[0:n-1]))

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}
		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		packet := &Packet{
			addr.Port,
			"here",
			4.4}
		err = encoder.Encode(packet)
		if err != nil {
			fmt.Println("--error")
			return
		}
		_, err = connection.WriteToUDP(buf.Bytes(), addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
