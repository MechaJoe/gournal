package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// for debugging
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type errMsg error

type model struct {
	textinput textinput.Model
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Today, I..."
	ti.Focus()

	return model{
		textinput: ti,
		err:       nil,
	}
}

func processEntry(text string) tea.Cmd {
	entry := Entry{
		Text: text,
		Time: time.Now(),
	}
	return func() tea.Msg {
		return errMsg(Save(entry))
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
		case tea.KeyEnter:
			entry := Entry{
				Text: m.textinput.Value(),
				Time: time.Now(),
			}
			Save(entry)
			// processEntry(m.textinput.Value())
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	// if it's the morning, show this prompt
	var message string

	currentHour := time.Now().Hour()
	if currentHour < 12 {
		message = "Good morning! It's currently %s. What do you want to accomplish today?"
	} else if currentHour < 17 {
		message = "Good afternoon. It's currently %s. How is your day going?"
	} else {
		message = "Good evening. It's currently %s. How was your day?"
	}
	formatted_time := fmt.Sprintf("%d:%d", currentHour, time.Now().Minute())
	message = fmt.Sprintf(message, formatted_time)

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		message,
		m.textinput.View(),
		"(ctrl+c to cancel, enter to record your journal entry)\n\n",
	)
}
