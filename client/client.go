package main

import (
	"bytes"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	conn    *websocket.Conn
	send    chan []byte
	recieve chan []byte
}

// writePump pumps recieved messages from server
func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Fatal(err)
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// readPump pumps messages a websocket
func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error: %s", err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		select {
		case c.recieve <- message:
		default:
			log.Println("Channel full, dropping message")
		}
	}
}

func NewClient(serverURL string) (*Client, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:    conn,
		send:    make(chan []byte, 256),
		recieve: make(chan []byte, 256),
	}
	return client, nil
}

func (c *Client) Start() {
	go c.writePump()
	go c.readPump()
}
func (c *Client) Send(message []byte) {
	c.send <- message
}
func (c *Client) Close() {
	close(c.send)
}
