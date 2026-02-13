package main

import "encoding/json"

type ClientMessage struct {
	clientMessageType int
	Content           string
}

type ServerMessage struct {
	Type    int
	Content string
}

func (c *Client) Join(nickname string) error {
	message := ClientMessage{
		clientMessageType: 0,
		Content:           nickname}

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return err
	}

	c.Send(jsonMsg)
	return nil
}
