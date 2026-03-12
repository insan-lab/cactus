// Copyright 2014 The Cactus Authors. All rights reserved.

package hub

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	conns = map[*websocket.Conn]bool{}

	chAdd = make(chan *websocket.Conn)
	chDel = make(chan *websocket.Conn)
	chMsg = make(chan interface{})
)

func Send(v interface{}) {
	chMsg <- v
}

func HandleConnect(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "", 400)
		return
	}
	catch(err)

	chAdd <- c

	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(5 * time.Second))
		return nil
	})

	go func() {
		defer func() {
			chDel <- c
			c.Close()
		}()

		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	go func() {
		defer func() {
			c.Close()
		}()

		for {
			<-time.After(3 * time.Second)

			err := c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		}
	}()

	err = c.WriteJSON([]interface{}{"HELO"})
	catch(err)
}

func init() {
	go func() {
		for {
			select {
			case c := <-chAdd:
				conns[c] = true

			case c := <-chDel:
				delete(conns, c)

			case v := <-chMsg:
				for c := range conns {
					if err := c.WriteJSON(v); err != nil {
						log.Printf("WebSocket write error: %v", err)
					}
				}
			}
		}
	}()
}
