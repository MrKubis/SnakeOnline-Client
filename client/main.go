package main

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type gameState int

const (
	stateConnecting gameState = iota
	statePlaying
	board_width  = 20
	board_height = 20
)

var serverUrl string = "ws://localhost:5250/ws"

//Board
// 0 - empty
// 1 - snake
// 2 - fruit

func parseBoardData(message string) []byte {
	result := strings.ReplaceAll(message, "\n", "")
	return []byte(result)
}

type Styles struct {
}

func DefaultStyles() *Styles {
	s := new(Styles)
	return s
}

func NewModel(client *Client) *model {
	styles := DefaultStyles()

	board := [board_height][board_width]int{}
	for i := range board_height {
		for j := range board_width {
			board[i][j] = 0
		}
	}

	return &model{
		state:  stateConnecting,
		styles: styles,
		client: client,
		board:  board,
	}
}

func (m *model) renderBoard() string {

	var result string = ""
	for i := range board_height {
		for j := range board_width {
			switch m.board[i][j] {
			case 0:
				result += "[ ]"
			case 1:
				result += "[X]"
			case 2:
				result += "[@]"
			}
		}
		result += "\n"
	}

	return result
}

type BoardUpdateMsg struct {
	board [board_height][board_width]int
}
type WebSocketMsg struct {
	data []byte
}

func main() {

	client, err := NewClient(serverUrl)
	if err != nil {
		log.Fatal("Error connection: ", err)
	}

	defer client.Close()

	client.Start()
	err = client.Join("MrOdbycik")

	if err != nil {
		log.Fatal(err)
	}

	m := NewModel(client)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
