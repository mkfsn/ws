package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Client interface {
}

type client struct {
	connection *websocket.Conn
	done       chan struct{}

	receiveCh chan []byte
}

func NewClient(url string, header http.Header) (*client, error) {
	log.Println("Dialing to", url, "...")
	connection, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		return nil, err
	}

	c := client{
		connection: connection,
		receiveCh:  make(chan []byte),
		done:       make(chan struct{}),
	}

	go c.receive()

	return &c, nil
}

func (c *client) Close() error {
	close(c.done)

	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := c.connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}

	return c.connection.Close()
}

func (c *client) Send(data []byte) error {
	err := c.connection.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Receive() <-chan []byte {
	return c.receiveCh
}

func (c *client) receive() {
	for {
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			log.Println("Failed to read:", err)
			return
		}
		c.receiveCh <- []byte(message)
	}
}
