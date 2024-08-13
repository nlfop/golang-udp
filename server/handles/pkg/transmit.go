package transmit

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"net"
	"time"
	"udp_connect/handles/structures"
)

func TransmitStructure(ctx context.Context, connection *net.UDPConn, addr *net.UDPAddr) {
	counter := 1
	for {
		timer1 := time.NewTimer(100 * time.Millisecond)
		select {
		case <-timer1.C:
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
				timer1.Stop()
				fmt.Println("--error")
				return
			}
			_, err = connection.WriteToUDP(buf.Bytes(), addr)
			if err != nil {
				timer1.Stop()
				fmt.Println(err)
				return
			}

		case <-ctx.Done():
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
