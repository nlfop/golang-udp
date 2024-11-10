package slave

//go run slave/slave.go 127.0.0.1:1234 127.0.0.1:1235 127.0.0.1:8000 127.0.0.1:8001

// 68000701000000 кс старт
// 68000702000000 кс стоп

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"
	"udp_connect/server_master/handles/command"
	transmit "udp_connect/server_master/handles/pkg"

	"github.com/gofiber/fiber/v2"
)

type reqString struct {
	CommandJSON string `json:"command"`
}

func ReceiveCommandFront(connectionhttp *fiber.Ctx) error {

	CONNECTComm := "127.0.0.1:1235"

	addrComm := "127.0.0.1:8001"

	sComm, err := net.ResolveUDPAddr("udp4", CONNECTComm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	sAddrComm, err := net.ResolveUDPAddr("udp4", addrComm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	cComm, err := net.DialUDP("udp4", sAddrComm, sComm)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Enter 7 bytes without spaces and checksum (it counts automatically)")

	defer cComm.Close()

	req := new(reqString)
	err = connectionhttp.BodyParser(req)
	fmt.Println(req.CommandJSON)
	data := req.CommandJSON
	if err != nil {

		connectionhttp.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
		return fmt.Errorf("something's wrong with your input")
	}

	str, _ := hex.DecodeString(data)
	str = transmit.CountCheckSum(str)
	fmt.Println(str)
	_, err = cComm.Write(str)
	if err != nil {
		fmt.Println("--error")
		return err
	}

	comm, err := command.CommandTrim(str)
	if err != nil {
		fmt.Println(err)
		connectionhttp.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
		return fmt.Errorf("something's wrong with your input")
	}
	switch comm {
	case "STOP_FLOW":

		time.Sleep(2000 * time.Millisecond)
		fmt.Println("Exiting UDP client flow data!")
		connectionhttp.Status(200).JSON(fiber.Map{"status": "Exiting UDP client!"})

	}
	return nil
}
