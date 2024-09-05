package transmit

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"time"
	"udp_connect/server/handles/structures"
	"unsafe"

	"github.com/gofiber/fiber/v2"
)

func ReceiveStructure(c *net.UDPConn, ctx context.Context) {

	nameFile := fmt.Sprintf("%v_%v.bin", time.Now().Format("2006_01_02"), time.Now().Format("15_04"))
	fileBIN, err := os.Create(nameFile)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer fileBIN.Close()
	nameFileTXT := fmt.Sprintf("%v_%v.txt", time.Now().Format("2006_01_02"), time.Now().Format("15_04"))
	fileTXT, err := os.Create(nameFileTXT)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}

	defer fileTXT.Close()
	for {

		select {
		case <-ctx.Done():
			fmt.Println("here")

			time.Sleep(200 * time.Millisecond)
			return
		default:
			buffer := make([]byte, 4096)

			c.SetDeadline(time.Now().Add(2 * time.Second))
			n, _, err := c.ReadFromUDP(buffer)
			if err != nil {
				continue
			}
			dec := gob.NewDecoder(bytes.NewReader(buffer[:n]))
			p := structures.Packet{}

			dec.Decode(&p)

			fileBIN.Write(buffer)
			fileTXT.WriteString(fmt.Sprintf("%d %v\n", p.PortTo, unsafe.Sizeof(p.PortTo)+unsafe.Sizeof(p.Message)+unsafe.Sizeof(p.NumFloat)+4*uintptr(len(p.BigMass))))

		}

	}
}

func ReceiveStructureOnce(connectionhttp *fiber.Ctx, c *net.UDPConn) {
	nameFile := fmt.Sprintf("%v_%v.bin", time.Now().Format("2006_01_02"), time.Now().Format("15_04"))
	fileBIN, err := os.Create(nameFile)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer fileBIN.Close()
	nameFileTXT := fmt.Sprintf("%v_%v.txt", time.Now().Format("2006_01_02"), time.Now().Format("15_04"))
	fileTXT, err := os.Create(nameFileTXT)
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}

	defer fileTXT.Close()
	buffer := make([]byte, 4096)

	n, _, _ := c.ReadFromUDP(buffer)
	fmt.Println(n)
	dec := gob.NewDecoder(bytes.NewReader(buffer[:n]))
	p := structures.Packet{}

	err = dec.Decode(&p)
	if err != nil {
		fmt.Println("Decode", err)
	}
	fileBIN.Write(buffer)
	fileTXT.WriteString(fmt.Sprintf("%d %v\n", p.PortTo, unsafe.Sizeof(p.PortTo)+unsafe.Sizeof(p.Message)+unsafe.Sizeof(p.NumFloat)+4*uintptr(len(p.BigMass))))

	connectionhttp.Status(200).JSON(p)

}

func CountCheckSum(b []byte) []byte {
	check := 0
	for i := 0; i < len(b); i++ {
		check += int(b[i])
	}
	check = check % 256
	b = append(b, byte(check))
	return b
}
