// Server-side part of the Go websocket sample.
//
// Eli Bendersky [http://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"flag"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// websocketTimeConnection handles a single websocket time connection - ws.
func websocketTimeConnection(ws *websocket.Conn) {
	for range time.Tick(1 * time.Second) {

		ws.WriteMessage(websocket.TextMessage, []byte(time.Now().Format(time.Stamp)))

	}
}

func main() {
	app := fiber.New()

	app.Get("/wstime", websocket.New(websocketTimeConnection))

	app.Static("/", "server/static/html/index.html")
	addr := flag.String("addr", ":8080", "http service address")
	flag.Parse()
	app.Listen(*addr)

}
