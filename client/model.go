package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	state        gameState
	isGameLoaded bool
	isGameOver   bool
	board        [board_height][board_width]int
	textInput    textinput.Model
	spinner      spinner.Model
	nickname     string
	width        int
	height       int
	styles       *Styles
	client       *Client
	keys         keyMap
	help         help.Model
}

type Styles struct {
	BorderColor lipgloss.Color
	SnakeStyle  lipgloss.Style
	FruitStyle  lipgloss.Style
	LoginStyle  lipgloss.Style
}

var (
	logo = `
 ▗▄▄▖▄▄▄▄  ▗▞▀▜▌█  ▄ ▗▞▀▚▖     ▗▄▖ ▗▖  ▗▖█ ▄ ▄▄▄▄  ▗▞▀▚▖
▐▌   █   █ ▝▚▄▟▌█▄▀  ▐▛▀▀▘    ▐▌ ▐▌▐▛▚▖▐▌█ ▄ █   █ ▐▛▀▀▘
 ▝▀▚▖█   █      █ ▀▄ ▝▚▄▄▖    ▐▌ ▐▌▐▌ ▝▜▌█ █ █   █ ▝▚▄▄▖
▗▄▄▞▘           █  █          ▝▚▄▞▘▐▌  ▐▌█ █            
`
)

const (
	stateConnecting gameState = iota
	statePlaying
	stateLogging
)

var (
	CustomSpinner = spinner.Spinner{
		Frames: []string{"🌍   Searching for opponent   🌍", "🌎 . Searching for opponent . 🌎", "🌏 .. Searching for opponent .. 🌏"},
		FPS:    time.Second / 2, //nolint:mnd
	}
)

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("#76bdff")
	s.SnakeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#76bdff"))
	s.FruitStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff3490"))
	s.LoginStyle = lipgloss.NewStyle()

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

func listenForWSMessages(msg <-chan []byte) tea.Cmd {
	return func() tea.Msg {
		return WebSocketMsg{data: <-msg}
	}
}

func parseBoardData(message string) []byte {
	result := strings.ReplaceAll(message, "\n", "")
	return []byte(result)
}

func (m *model) renderBoard() string {

	var result string = ""
	for i := range board_height {
		for j := range board_width {
			switch m.board[i][j] {
			case 0:
				result += "  "
			case 1:
				result += m.styles.SnakeStyle.Render("██")
			case 2:
				result += m.styles.FruitStyle.Render("██")

			}
		}
		if i < board_height-1 {
			result += "\n"
		}
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.styles.BorderColor)

	helpview := m.help.View(m.keys)

	if !m.isGameOver {
		return lipgloss.JoinVertical(lipgloss.Center, style.Render(result), helpview)
	} else {
		return lipgloss.JoinVertical(lipgloss.Center, style.Render(result))
	}
}

func (m *model) renderLogging() string {

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.styles.BorderColor).
		Padding(1, 2)

	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#76bdff"))
	input := style.Render(lipgloss.JoinVertical(lipgloss.Center, "Enter your nickname\n", m.styles.LoginStyle.Render(m.textInput.View())))
	return lipgloss.JoinVertical(lipgloss.Center, logoStyle.Render(logo), "\n\n", input)
}

func (m *model) renderLoading() string {
	return m.spinner.View()
}

type BoardUpdateMsg struct {
	board [board_height][board_width]int
}
type WebSocketMsg struct {
	data []byte
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
func (m *model) resetBoard() {
	for i := range board_height {
		for j := range board_width {
			m.board[i][j] = 0
		}
	}
}

func (m *model) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Spinner = CustomSpinner
}

func initialModel(c *Client) model {
	styles := DefaultStyles()
	board := [board_height][board_width]int{}
	for i := range board_height {
		for j := range board_width {
			board[i][j] = 0
		}
	}

	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	sp := spinner.New()
	sp.Spinner = CustomSpinner

	return model{
		state:     stateLogging,
		styles:    styles,
		client:    c,
		board:     board,
		textInput: ti,
		spinner:   sp,

		keys: keys,
		help: help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		if m.state == stateLogging {
			switch {

			case msg.String() == "ctrl+c":
				return m, tea.Quit

			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll

			case key.Matches(msg, m.keys.Enter):
				if m.textInput.Value() != "" {
					m.nickname = m.textInput.Value()
					err := m.client.Join(m.nickname)
					if err != nil {
						log.Fatal(err)
						return m, nil
					}
					m.state = statePlaying
					m.resetSpinner()
					return m, tea.Batch(
						listenForWSMessages(m.client.recieve),
						m.spinner.Tick,
					)
				}

			}
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
		if m.state == statePlaying {

			switch {

			case key.Matches(msg, m.keys.Help):
				m.help.ShowAll = !m.help.ShowAll

			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.keys.Up):
				m.client.Move("Up")

			case key.Matches(msg, m.keys.Down):
				m.client.Move("Down")

			case key.Matches(msg, m.keys.Left):
				m.client.Move("Left")

			case key.Matches(msg, m.keys.Right):
				m.client.Move("Right")
			case key.Matches(msg, m.keys.Enter):
				if m.isGameOver {
					m.state = stateLogging
					m.isGameLoaded = false
					m.isGameOver = false
					m.resetBoard()
					return m, tea.Batch(textinput.Blink, m.spinner.Tick)
				}

			}
		}
	case WebSocketMsg:

		serverMsg := &ServerMessage{}
		err := json.Unmarshal(msg.data, serverMsg)

		if err != nil {
			return m, cmd
		}
		switch serverMsg.Type {
		case 4:
			m.isGameOver = true

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

func (m model) View() string {

	result := ""
	switch m.state {

	case stateLogging:
		result += m.renderLogging()
	case stateConnecting:
		result += m.renderLoading()
	case statePlaying:
		if !m.isGameLoaded {
			result += m.renderLoading() + "\n"
		} else {
			result += "\n"
		}
		result += m.renderBoard()
	}

	overlay := ""

	if m.isGameOver {

		overlay = "Game over!\nPress <ENTER> to quit"
		overlay = lipgloss.Place(
			m.width,
			2,
			lipgloss.Center,
			lipgloss.Top,
			overlay,
		)
	}
	centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, result, overlay))

	return centered
}
