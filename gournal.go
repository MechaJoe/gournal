package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

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
	InitJournal()
	var p *tea.Program

	journal, _ := GetJournal()
	if journal == nil || journal.User == "" {
		p = tea.NewProgram(initialModel())
	} else if len(os.Args) < 2 || os.Args[1] == "new" {
		p = tea.NewProgram(writingModel())
	} else if len(os.Args) == 2 && os.Args[1] == "browse" {
		p = tea.NewProgram(browsingModel())
	}

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type errMsg error

type initModel struct {
	textinput textinput.Model
	err       error
}

type journalModel struct {
	textinput textinput.Model
	err       error
}

type browseModel struct {
	list list.Model
	err  error
}

func initialModel() initModel {
	ti := textinput.New()
	ti.Placeholder = "Adam"
	ti.Focus()

	return initModel{
		textinput: ti,
		err:       nil,
	}
}

func browsingModel() browseModel {
	journal, err := GetJournal()
	if err != nil {
		log.Fatal(err)
	}
	items := make([]list.Item, len(journal.Entries))
	for i := 0; i < len(journal.Entries); i++ {
		items[i] = journal.Entries[i]
	}

	m := browseModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0), err: nil}
	m.list.Title = fmt.Sprintf("%s's Gournal Entries", journal.User)
	return m
}

func writingModel() journalModel {
	ti := textinput.New()
	ti.Placeholder = "Today, I..."
	ti.Focus()

	return journalModel{
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

func (m initModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m journalModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m browseModel) Init() tea.Cmd {
	return nil
}

func (m initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			UpdateName(m.textinput.Value())
			
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m browseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m journalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m browseModel) View() string {
	return docStyle.Render(m.list.View())
}

func (m initModel) View() string {
	greetingMessage := "Greetings! Let's get started with gournal--what's your name?"
	salutation := "After you're done, run \"gournal new\" to record your first gournal entry!"
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n",
		greetingMessage,
		m.textinput.View(),
		salutation,
	)
}

func (m journalModel) View() string {
	var message string
	journal, _ := GetJournal()
	name := journal.User
	currentHour := time.Now().Hour()
	if currentHour < 12 {
		message = "Good morning, %s! It's currently %s. What do you want to accomplish today?"
	} else if currentHour < 17 {
		message = "Good afternoon, %s. It's currently %s. How is your day going?"
	} else {
		message = "Good evening, %s. It's currently %s. How was your day?"
	}
	formatted_time := fmt.Sprintf("%02d:%02d", currentHour, time.Now().Minute())
	message = fmt.Sprintf(message, name, formatted_time)

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n",
		message,
		m.textinput.View(),
		"(ctrl+c to cancel, enter to record your journal entry)",
	)
}
