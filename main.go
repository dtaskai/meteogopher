package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const geocodingApiUrl = "https://geocoding-api.open-meteo.com/v1/search?name=Berlin"
const forecastApiUrl = "https://api.open-meteo.com/v1/forecast"

func main() {
	p := tea.NewProgram(initialModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

type statusMsg int
type errMsg struct{ err error }
type model struct {
	textInput textinput.Model
	err       error
	data      int
}

func initialModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			msg := checkServer()
			switch msg := msg.(type) {
			case statusMsg:
				m.data = int(msg)
				return m, nil
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := fmt.Sprintf(
		"Enter a location:\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"

	if m.err != nil {
		s += fmt.Sprintf("something went wrong: %s", m.err)
	} else if m.data != 0 {
		s = fmt.Sprintf("%d %s", m.data, http.StatusText(m.data))
	}

	return s + "\n"
}

func checkServer() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(geocodingApiUrl)

	if err != nil {
		return errMsg{err}
	}

	return statusMsg(res.StatusCode)
}

func (e errMsg) Error() string { return e.err.Error() }
