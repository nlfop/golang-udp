package main

//go run slave/slave.go 127.0.0.1:1234

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

type Packet struct {
	PortTo   int
	Message  string
	NumFloat float32
}

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
	defer c.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text + "\n")
		_, err = c.Write(data)
		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Println("Exiting UDP client!")
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}
		buffer := make([]byte, 4096)
		n, _, _ := c.ReadFromUDP(buffer)

		dec := gob.NewDecoder(bytes.NewReader(buffer[:n]))
		p := Packet{}
		dec.Decode(&p)

		// n, _, err := c.ReadFromUDP(buffer)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		fmt.Println(p)
	}
}
