package main

import (
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	state        gameState
	isGameLoaded bool
	board        [board_height][board_width]int
	width        int
	height       int
	styles       *Styles
	client       *Client
}

func listenForWSMessages(msg <-chan []byte) tea.Cmd {
	return func() tea.Msg {
		return WebSocketMsg{data: <-msg}
	}
}

func (m model) Init() tea.Cmd {
	return listenForWSMessages(m.client.recieve)
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.state == statePlaying {
				m.client.Move("Up")
			}
		case "down":
			if m.state == statePlaying {
				m.client.Move("Down")
			}
		case "left":
			if m.state == statePlaying {
				m.client.Move("Left")
			}
		case "right":
			if m.state == statePlaying {
				m.client.Move("Right")
			}
		}
	case WebSocketMsg:

		serverMsg := &ServerMessage{}
		err := json.Unmarshal(msg.data, serverMsg)

		if err != nil {
			return m, cmd
		}

		switch serverMsg.Type {
		case 5:
			boardData := parseBoardData(serverMsg.Content)
			m.UpdateBoard(boardData)
			m.isGameLoaded = true
			m.state = statePlaying
		}

		return m, listenForWSMessages(m.client.recieve)
	}
	return m, cmd
}

func (m *model) UpdateBoard(data []byte) {

	index := 0
	for i := range board_height {
		for j := range board_width {
			if index >= len(data) {
				return
			}
			switch string(data[index]) {
			case "0":
				m.board[i][j] = 0
			case "S":
				m.board[i][j] = 1
			case "F":
				m.board[i][j] = 2
			}
			index++
		}
	}
}

func (m model) View() string {
	switch m.state {
	case stateConnecting:
		return "Loading..."
	case statePlaying:
		if m.width == 0 || !m.isGameLoaded {
			return "loading..."
		}
		return m.renderBoard()
	}
	return ""
}
