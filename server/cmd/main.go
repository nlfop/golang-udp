// Server-side part of the Go websocket sample.
//
// Eli Bendersky [http://eli.thegreenplace.net]
// This code is in the public domain.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

var (
	port = flag.Int("port", 4050, "The server port")
)

// websocketTimeConnection handles a single websocket time connection - ws.
func websocketTimeConnection(ws *websocket.Conn) {
	for range time.Tick(1 * time.Second) {
		// Once a second, send a message (as a string) with the current time.
		websocket.Message.Send(ws, time.Now().Format(time.StampMilli))
	}
}

func main() {
	flag.Parse()
	// Set up websocket servers and static file server. In addition, we're using
	// net/trace for debugging - it will be available at /debug/requests.

	http.Handle("/wstime", websocket.Handler(websocketTimeConnection))
	http.Handle("/", http.FileServer(http.Dir("server/static/html")))

	log.Printf("Server listening on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
