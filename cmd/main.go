package main

import (
	"app/internal/tui"
	"log"

	"github.com/MrKubis/SnakeOnline-Client/internal/client"
	"github.com/MrKubis/SnakeOnline-Client/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type gameState int

const (
	board_width  = 20
	board_height = 20
)

var serverUrl string = "ws://localhost:5250/ws"

//Board
// 0 - empty
// 1 - snake
// 2 - fruit

func main() {

	client, err := client.NewClient(serverUrl)
	if err != nil {
		log.Fatal("Error connection: ", err)
	}

	defer client.Close()

	client.Start()

	p := tea.NewProgram(tui.InitialModel(client), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
