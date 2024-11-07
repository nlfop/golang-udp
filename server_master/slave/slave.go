package slave

//go run slave/slave.go 127.0.0.1:1234 127.0.0.1:1235 127.0.0.1:8000 127.0.0.1:8001

// 68 00 10 10 00 00 00 - старт периодической передачи
// 68 00 20 01 00 00 00 - стоп периодической передачи

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
	// arguments := os.Args
	// if len(arguments) == 1 {
	// 	fmt.Println("Please provide a host:port string")
	// 	return
	// }

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
	// var ctx context.Context
	// var cancel context.CancelFunc

	req := new(reqString) // Store the body in the user and return error if encountered
	err = connectionhttp.BodyParser(req)
	fmt.Println(req.CommandJSON)
	data := req.CommandJSON
	if err != nil {
		// cancel()
		connectionhttp.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
		return fmt.Errorf("something's wrong with your input")
	}
	// Return the created user
	// n := len(data)

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

		// case "START_FLOW":
		// 	ctx, cancel = context.WithCancel(context.Background())

		// 	go transmit.ReceiveStructure(cData, ctx, cancel)
		// case "START_ONCE":
		// 	go transmit.ReceiveStructureOnce(cData)
	}
	return nil
}
