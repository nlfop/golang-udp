package transmit

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"time"
	"udp_connect/command_rt/handles/command"
	"udp_connect/command_rt/handles/structures"
)

func TransmitStructure(ctx context.Context, cancel context.CancelFunc, connection *net.UDPConn, addr *net.UDPAddr) {
	file, err := os.Open("Crate.bin")
	if err != nil {

		cancel()
		panic(err)

	}
	defer file.Close()
	numBytes := 8
	sizePacket := -1
	var FirstPackage bool
	for {
		timer1 := time.NewTimer(100 * time.Millisecond)
		select {
		case <-timer1.C:
			data := make([]byte, numBytes)
			n, err := file.Read(data)
			if err != nil || n == 0 { // если конец файла
				file.Seek(0, 0)
				continue
			}

			numBytes = CountPackage(data)

			_, err = connection.WriteToUDP(data, addr)
			if err != nil {
				timer1.Stop()
				fmt.Println(err)
				return
			}
			if sizePacket > 0 {
				sizePacket -= n
				if sizePacket == 0 {
					fmt.Println("packet!")
					if !FirstPackage {
						answer := []byte{0x53, 0x00, 0x07, 0x00, 0x00, 0x01, 0x00}
						answer = CountCheckSum(answer)
						command.CommandTrim(answer)
						connection.WriteToUDP(answer, addr)
						FirstPackage = true
					}
				}
			}

			if data[0] == 83 {
				sl := []byte{data[2], data[1]}
				dataSize := binary.BigEndian.Uint16(sl)
				sizePacket = int(dataSize)
			}
			// fmt.Printf("%x\n", data)
		case <-ctx.Done():
			for sizePacket > 0 {
				data := make([]byte, numBytes)
				n, err := file.Read(data)
				if err != nil || n == 0 { // если конец файла
					file.Seek(0, 0)
					continue
				}

				numBytes = CountPackage(data)

				_, err = connection.WriteToUDP(data, addr)
				if err != nil {
					timer1.Stop()
					fmt.Println(err)
					return
				}

				sizePacket -= n
				if sizePacket == 0 {
					fmt.Println("packet!")
					if !FirstPackage {
						answer := []byte{0x53, 0x00, 0x07, 0x00, 0x00, 0x01, 0x00}
						answer = CountCheckSum(answer)
						command.CommandTrim(answer)
						connection.WriteToUDP(answer, addr)
						FirstPackage = true
					}
				}
			}
			answer := []byte{0x53, 0x00, 0x07, 0x00, 0x00, 0x02, 0x00}
			answer = CountCheckSum(answer)
			command.CommandTrim(answer)
			connection.WriteToUDP(answer, addr)

			timer1.Stop()

			return
		}

	}
}

func TransmitStructureOnce(connection *net.UDPConn, addr *net.UDPAddr) {
	counter := 1

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	packet := &structures.Packet{
		PortTo:   counter,
		Message:  "here",
		NumFloat: 4.4,
		BigMass:  make([]int32, 256)}
	counter += 1

	err := encoder.Encode(packet)
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

func CountPackage(d []byte) int {
	switch {
	case len(d) != 8:

		return 8
	case d[0] == 83:

		return 8

	default:

		if d[0] == 0 {

			return 8
		}

		return int(d[0])
	}

}
