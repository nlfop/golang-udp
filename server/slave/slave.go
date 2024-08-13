package slave

//go run slave/slave.go 127.0.0.1:1234

import (
	"context"
	"encoding/hex"
	"fmt"
	"net"
	"time"
	"udp_connect/handles/command"

	transmit "udp_connect/handles/pkg"

	"github.com/gofiber/fiber/v2"
)

type reqString struct {
	CommandJSON string `json:"command"`
}

func ReceiveCommandFront(connectionhttp *fiber.Ctx) error {
	CONNECT := "127.0.0.1:1234"
	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	fmt.Println("Enter 7 bytes without spaces and checksum (it counts automatically)")
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	// for {

	// n := len(data)
	req := new(reqString) // Store the body in the user and return error if encountered
	err = connectionhttp.BodyParser(req)
	fmt.Println(req.CommandJSON)
	data := req.CommandJSON
	if err != nil {
		return connectionhttp.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}
	// Return the created user
	n := len(data)

	str, _ := hex.DecodeString(data[0 : n-1])
	str = transmit.CountCheckSum(str)
	_, err = c.Write(str)
	if err != nil {
		fmt.Println("--error")
		return nil
	}
	fmt.Println(str)
	comm, err := command.CommandTrim(str)
	if err != nil {
		fmt.Println(err)
		return connectionhttp.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}
	switch comm {
	case "STOP_FLOW":
		cancel()
		time.Sleep(2000 * time.Millisecond)
		fmt.Println("Exiting UDP client!")
		return nil
	case "START_FLOW":
		go transmit.ReceiveStructure(c, ctx)
	case "START_ONCE":
		go transmit.ReceiveStructureOnce(c)
	}

	return connectionhttp.Status(200).JSON(fiber.Map{"status": "success"})
}
