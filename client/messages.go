package main

import (
	"encoding/json"
	"fmt"
)

type ClientMessage struct {
	Type    string
	Content string
}

type ServerMessage struct {
	Type    int
	Content string
}

func (c *Client) Join(nickname string) error {
	message := ClientMessage{
		Type:    "JOIN",
		Content: nickname}

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	c.Send(jsonMsg)
	return nil
}

func (c *Client) Move(direction string) error {

	fmt.Println(direction)
	message := ClientMessage{
		Type:    "MOVE",
		Content: direction,
	}

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	c.Send(jsonMsg)
	return nil
}
