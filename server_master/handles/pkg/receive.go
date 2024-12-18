package transmit

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
	"udp_connect/server_master/handles/command"

	"github.com/gofiber/contrib/websocket"
)

type ReciveStructure struct {
	Prefix      byte                    `json:"prefix"`
	Size        int                     `json:"size"`
	Command     byte                    `json:"command"`
	ErrorCode   byte                    `json:"error_code"`
	Info1       byte                    `json:"info_1"`
	Info2       byte                    `json:"info_2"`
	ControlSumm byte                    `json:"control_summ"`
	Data        []ReceiveStructureBlock `json:"data"`
}

type ReceiveStructureBlock struct {
	Size                byte   `json:"size"`
	ErrorBlock          byte   `json:"error_block"`
	TypeStructure       byte   `json:"type_str"`
	NumberPlace         byte   `json:"num_block"`
	Info2               byte   `json:"info_2"`
	Info3               byte   `json:"info_3"`
	ControlSummSructure byte   `json:"control_str"`
	ControlSumm         byte   `json:"control_summ"`
	DataBlock           string `json:"data_block"`
}

type CommandAnswer struct {
	Answer string `json:"commandAnswer"`
}

func ReceiveStructure(c *net.UDPConn, ctx context.Context, cancel context.CancelFunc, ws *websocket.Conn) {
	answer = &ReciveStructure{}
	block = &ReceiveStructureBlock{}
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

	buffer := make([]byte, 4096)
	sizePacket := -1
	var FirstPackage int

	defer fileTXT.Close()
	for {

		select {
		case <-ctx.Done():

			time.Sleep(200 * time.Millisecond)

			return
		default:

			c.SetDeadline(time.Now().Add(2 * time.Second))
			n, _, err := c.ReadFromUDP(buffer)
			if err != nil {
				cancel()
				fmt.Println(err)
				continue
			}
			if n == 8 {
				if buffer[0] == 83 && buffer[5] == 2 {
					fileTXT.WriteString(fmt.Sprintf(">> Ответ на команду: %x\n", buffer[0:n]))

					answers = append(answers, *answer)

					dataAns, _ := json.Marshal(answer)
					ws.WriteMessage(websocket.TextMessage, dataAns)
					jsonBuf := CommandAnswer{
						Answer: fmt.Sprintf("%x", buffer[0:n]),
					}
					data, _ := json.Marshal(jsonBuf)

					ws.WriteMessage(websocket.TextMessage, data)
					_, err := command.CommandTrim(buffer[0:n])
					if err != nil {
						fmt.Println(err)
					}
					cancel()
					continue
				}
			}

			fileBIN.Write(buffer[:n])

			if buffer[0] == 83 && !commWait {
				sl := []byte{buffer[2], buffer[1]}
				dataSize := binary.BigEndian.Uint16(sl)
				sizePacket = int(dataSize)
			}
			EncodingPackage(buffer[:n], fileTXT, ws)
			if sizePacket > 0 {
				sizePacket -= n
				if sizePacket == 0 {
					FirstPackage++
				}
			}

		}

	}
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

var answers []ReciveStructure
var answer = &ReciveStructure{}
var block = &ReceiveStructureBlock{}

var commWait bool

func EncodingPackage(d []byte, fileTXT *os.File, ws *websocket.Conn) {
	switch {
	case len(d) != 8 || commWait:
		block.DataBlock = fmt.Sprintf("%x", d)
		answer.Data = append(answer.Data, *block)
		fileTXT.WriteString(fmt.Sprintf("		Данные : %x\n", d))
		commWait = false
	case d[0] == 83:
		if answer.Prefix == 83 {
			answers = append(answers, *answer)
			data, _ := json.Marshal(answer)

			ws.WriteMessage(websocket.TextMessage, data)
		}

		sl := []byte{d[2], d[1]}
		dataSize := binary.BigEndian.Uint16(sl)

		answer = &ReciveStructure{
			Prefix:      d[0],
			Size:        int(dataSize),
			Command:     d[3],
			ErrorCode:   d[4],
			Info1:       d[5],
			Info2:       d[6],
			ControlSumm: d[7],
		}
		fileTXT.WriteString(fmt.Sprintf("Новый пакет, размер %d\n", answer.Size))
	default:
		block = &ReceiveStructureBlock{
			Size:                d[0],
			ErrorBlock:          d[1],
			TypeStructure:       d[2],
			NumberPlace:         d[3],
			Info2:               d[4],
			Info3:               d[5],
			ControlSummSructure: d[6],
			ControlSumm:         d[7],
		}
		if block.TypeStructure == 2 {
			fileTXT.WriteString(fmt.Sprintf("	Данные команды, размер структуры %d\n", block.Size))
			commWait = true
		} else {
			fileTXT.WriteString(fmt.Sprintf("	Данные блока: %d, размер структуры %d\n", block.NumberPlace, block.Size))
		}

		if d[0] == 0 {
			answer.Data = append(answer.Data, *block)

		}

	}

}
