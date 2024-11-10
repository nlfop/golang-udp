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

func websocketTimeConnection(ws *websocket.Conn) {
	for range time.Tick(1 * time.Second) {

		ws.WriteMessage(websocket.TextMessage, []byte(time.Now().Format(time.Stamp)))

	}
}

func websocketFlowSlave(ws *websocket.Conn) {
	ctx, cancel := context.WithCancel(context.Background())
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
}

func main() {
	app := fiber.New()

	app.Get("/wstime", websocket.New(websocketTimeConnection))
	app.Get("/wsstruct", websocket.New(websocketStructureSend))
	app.Get("/flowSlave", websocket.New(websocketFlowSlave))
	app.Static("/", "././static/html/index.html")
	app.Post("/command", slave.ReceiveCommandFront)
	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()
	app.Listen(*addr)

}
