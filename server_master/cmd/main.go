// Server-side part of the Go websocket sample.
//
// Eli Bendersky [http://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"
	transmit "udp_connect/server_master/handles/pkg"
	"udp_connect/server_master/slave"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var messages = []string{"sky", "river", "mountain", "forest", "stone", "light", "ocean", "dream", "echo", "whisper", "star", "cloud", "fire", "ice", "shadow", "wind", "earth", "flame", "sun", "moon"}
var count int
var numStruct int

type ExampleStruct struct {
	Num     int    `json:"num"`
	Message string `json:"message"`
}

func countMessage() string {
	count++
	if count >= len(messages) {
		count = 0
	}
	return messages[count]
}

func newStruct() ExampleStruct {
	numStruct++
	return ExampleStruct{
		Num:     numStruct,
		Message: countMessage(),
	}
}

func websocketStructureSend(ws *websocket.Conn) {

	for range time.Tick(1 * time.Second) {
		structJSON := newStruct()
		data, _ := json.Marshal(structJSON)

		ws.WriteMessage(websocket.TextMessage, data)
	}
}

// websocketTimeConnection handles a single websocket time connection - ws.
func websocketTimeConnection(ws *websocket.Conn) {
	for range time.Tick(1 * time.Second) {

		ws.WriteMessage(websocket.TextMessage, []byte(time.Now().Format(time.Stamp)))

	}
}

func websocketFlowSlave(ws *websocket.Conn) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	CONNECTData := "127.0.0.1:1234"

	addrData := "127.0.0.1:8000"

	sData, err := net.ResolveUDPAddr("udp4", CONNECTData)
	if err != nil {
		fmt.Println(err)
		return
	}

	sAddrData, err := net.ResolveUDPAddr("udp4", addrData)
	if err != nil {
		fmt.Println(err)
		return
	}

	cData, err := net.DialUDP("udp4", sAddrData, sData)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cData.Close()
	transmit.ReceiveStructure(cData, ctx, cancel, ws)
	// CONNECT := "127.0.0.1:1234"
	// s, err := net.ResolveUDPAddr("udp4", CONNECT)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// c, err := net.DialUDP("udp4", nil, s)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	// // fmt.Println("Enter 7 bytes without spaces and checksum (it counts automatically)")
	// defer c.Close()

	// _, err = c.Write([]byte{104, 0, 16, 16, 0, 0, 0, 136})
	// nameFile := fmt.Sprintf("%v_%v.bin", time.Now().Format("2006_01_02"), time.Now().Format("15_04"))
	// fileBIN, err := os.Create(nameFile)
	// if err != nil {
	// 	fmt.Println("Unable to create file:", err)
	// 	os.Exit(1)
	// }
	// defer fileBIN.Close()
	// nameFileTXT := fmt.Sprintf("%v_%v.txt", time.Now().Format("2006_01_02"), time.Now().Format("15_04"))
	// fileTXT, err := os.Create(nameFileTXT)
	// if err != nil {
	// 	fmt.Println("Unable to create file:", err)
	// 	os.Exit(1)
	// }

	// defer fileTXT.Close()
	// for {

	// 	// select {
	// 	// case <-ctx.Done():
	// 	// 	fmt.Println("here")

	// 	// 	time.Sleep(200 * time.Millisecond)
	// 	// 	return
	// 	// default:
	// 	buffer := make([]byte, 4096)

	// 	c.SetDeadline(time.Now().Add(2 * time.Second))
	// 	n, _, err := c.ReadFromUDP(buffer)
	// 	if n == 0 {
	// 		fmt.Println("end of flow receive")
	// 		return
	// 	}
	// 	if err != nil {
	// 		continue
	// 	}
	// 	dec := gob.NewDecoder(bytes.NewReader(buffer[:n]))
	// 	p := structures.Packet{}

	// 	dec.Decode(&p)

	// 	fileBIN.Write(buffer)
	// 	fileTXT.WriteString(fmt.Sprintf("%d %v\n", p.PortTo, unsafe.Sizeof(p.PortTo)+unsafe.Sizeof(p.Message)+unsafe.Sizeof(p.NumFloat)+4*uintptr(len(p.BigMass))))

	// 	data, _ := json.Marshal(p)

	// ws.WriteMessage(websocket.TextMessage, data)
	// 	// }

	// }
}

func main() {
	app := fiber.New()

	app.Get("/wstime", websocket.New(websocketTimeConnection))
	app.Get("/wsstruct", websocket.New(websocketStructureSend))
	app.Get("/flowSlave", websocket.New(websocketFlowSlave))
	// app.Static("/", "udp_connect/server_master/static/html/index.html")

	app.Static("/", "././static/html/index.html")
	app.Post("/command", slave.ReceiveCommandFront)
	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()
	app.Listen(*addr)

}
