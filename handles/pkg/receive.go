package transmit

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
	"udp_connect/handles/structures"
	"unsafe"
)

func ReadSTOP(cancel context.CancelFunc, c *net.UDPConn) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">>> ")
	text, _ := reader.ReadString('\n')
	data := []byte(text + "\n")
	c.Write(data)
	if strings.TrimSpace(string(data)) == "STOP" {
		cancel()
		time.Sleep(2000 * time.Millisecond)
		fmt.Println("Exiting UDP client!")

	}
}

func ReceiveStructure(c *net.UDPConn, wg *sync.WaitGroup) {
	ctx, cancel := context.WithCancel(context.Background())
	go ReadSTOP(cancel, c)
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
			wg.Done()
			time.Sleep(200 * time.Millisecond)
			return
		default:
			buffer := make([]byte, 4096)

			c.SetDeadline(time.Now().Add(2 * time.Second))
			n, _, _ := c.ReadFromUDP(buffer)
			fileBIN.Write(buffer)
			dec := gob.NewDecoder(bytes.NewReader(buffer[:n]))
			p := structures.Packet{}
			dec.Decode(&p)

			fileTXT.WriteString(fmt.Sprintf("%d %v\n", p.PortTo, unsafe.Sizeof(p.PortTo)+unsafe.Sizeof(p.Message)+unsafe.Sizeof(p.NumFloat)+4*uintptr(len(p.BigMass))))

		}

	}
}
