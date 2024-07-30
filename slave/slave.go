package main

//go run slave/slave.go 127.0.0.1:1234

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
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
	defer c.Close()
	var wg sync.WaitGroup
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		data := []byte(text + "\n")
		_, err = c.Write(data)

		if strings.TrimSpace(string(data)) == "START" {
			wg.Add(1)
			go transmit.ReceiveStructure(c, &wg)
			wg.Wait()
			return
		}

	}
}
