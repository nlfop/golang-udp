// Server-side part of the Go websocket sample.
//
// Eli Bendersky [http://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"encoding/json"
	"flag"
	"time"
	"udp_connect/server/slave"

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

func main() {
	app := fiber.New()

	app.Get("/wstime", websocket.New(websocketTimeConnection))
	app.Get("/wsstruct", websocket.New(websocketStructureSend))
	app.Static("/", "server/static/html/index.html")
	app.Post("/command", slave.ReceiveCommandFront)
	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()
	app.Listen(*addr)

}
