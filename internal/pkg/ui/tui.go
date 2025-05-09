package ui

import (
	"fmt"
	"log"
	"os"
	//"github.com/charmbracelet/bubbles/list"
	//"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Application state type
type AppState int

// Enum for managing app state
const (
	tasksView = iota
	logView
)

// Styles
type Styles struct {
	BorderColor lipgloss.Color
	BorderStyle lipgloss.Style
	BoldText    lipgloss.Style
	Underlined  lipgloss.Style
}

func DefaultStyle(width int) *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("#e28743")
	s.BorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(s.BorderColor).Padding(0, 1)
	s.BoldText = lipgloss.NewStyle().Bold(true)
	s.Underlined = lipgloss.NewStyle().Underline(true)
	return s
}

// Root Model
type RootModel struct {
	state  AppState
	tasks  tea.Model
	log    tea.Model
	width  int
	height int
	style  *Styles
	data   map[string]interface{}
}

func (r RootModel) Init() tea.Cmd {
	return textinput.Blink
}

func (r RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return r, tea.Quit
		}
	}

	return r, cmd
}

func (r RootModel) View() string {
	title, ok := r.data["title"].(string)
	comment, _ := r.data["comment"].(string)
	if !ok {
		title = "No job title ..."
	}
	header := r.style.
		BorderStyle.Width(r.width-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(r.style.BoldText.Render("🚜 GO Tractor ! 🚜"))
	jobData := r.style.BorderStyle.Width(r.width-2).
		Align(lipgloss.Left, lipgloss.Center).
		Render(
			r.style.Underlined.Render("Job title:"),
			title,
			r.style.Underlined.Render("\nComment   :"),
			comment,
		)
	return lipgloss.Place(r.width, r.height, lipgloss.Center, lipgloss.Top, lipgloss.JoinVertical(lipgloss.Top, header, jobData))
}

func Show(data map[string]interface{}) {

	main := &RootModel{data: data}
	style := DefaultStyle(main.width)
	main.style = style

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	f.WriteString("Hello")
	p := tea.NewProgram(*main, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
