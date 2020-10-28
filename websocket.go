package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// We'll need to define an UpGrader
// this will require a Read and Write buffer size
var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// upgrade this connection to a WebSocket
	// connection
	ws, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// say hello
	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))

	// handle writer
	go func() {
		for {
			time.Sleep(time.Second * 5)
			err = ws.WriteMessage(1, []byte(time.Now().Format("2006-01-02 15:04:05")))
			if err != nil {
				log.Println(err)
			}
		}
	}()

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	go func(conn *websocket.Conn) {
		for {
			// read in a message
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			// print out that message for clarity
			log.Println(string(p))

			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}
		}
	}(ws)
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("starting server...")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
