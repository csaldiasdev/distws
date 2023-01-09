package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:60000", "http service address")
var token = flag.String("t", "", "user token")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	messageChan := make(chan []byte)
	disconectChan := make(chan interface{})

	query := url.Values{}
	query.Add("token", *token)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws", RawQuery: query.Encode()}

	firstTry := true

	for {
		if firstTry {
			log.Printf("connecting to %s", u.String())
		} else {
			log.Printf("reconnecting to %s", u.String())
		}

		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial:", err)
		}

		go func() {
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					c.Close()
					disconectChan <- true
					return
				}
				messageChan <- message
			}
		}()

	loop:
		for {
			select {
			case msg := <-messageChan:
				log.Printf(">>>> : %s", msg)
			case <-disconectChan:
				firstTry = false
				break loop
			case <-interrupt:
				c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}

		}
	}
}
