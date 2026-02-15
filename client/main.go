package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type gameState int

const (
	stateConnecting gameState = iota
	statePlaying
	stateLogging
	board_width  = 20
	board_height = 20
)

var serverUrl string = "ws://localhost:5250/ws"

//Board
// 0 - empty
// 1 - snake
// 2 - fruit

func main() {

	client, err := NewClient(serverUrl)
	if err != nil {
		log.Fatal("Error connection: ", err)
	}

	defer client.Close()

	client.Start()

	p := tea.NewProgram(initialModel(client), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
