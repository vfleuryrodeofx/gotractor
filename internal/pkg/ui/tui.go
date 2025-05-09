// // TODO :
// -
package ui

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/victorfleury/gotractor/internal/pkg/requests"
	"github.com/victorfleury/gotractor/internal/pkg/utils"

	//"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
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

// STYLES TO FIX
// Define some basic styling
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 0).
			Width(100).
			Align(lipgloss.Center) //.
		//Border(lipgloss.RoundedBorder())

	sectionStyle      = lipgloss.NewStyle().Padding(1, 2).Width(100)
	containerStyle    = lipgloss.NewStyle().Padding(1, 2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

// LIST Widget
type item string
type itemDelegate struct{}

func (d itemDelegate) Height() int { return 1 }

func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	var str string
	if index < 10 {
		str = fmt.Sprintf("%d.  %s", index+1, i)
	} else {
		str = fmt.Sprintf("%d. %s", index+1, i)
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
func (i item) FilterValue() string { return "" }

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
	state       AppState
	tasks       list.Model
	logViewport viewport.Model
	logContent  string
	width       int
	height      int
	style       *Styles
	data        map[string]any
	//tasksData []map[string]any
	tasksData []any
	jid       string
}

func initModel(data map[string]any, tasksData []any, jid string) *RootModel {

	//fmt.Println(data)

	// Initialize the list
	items := []list.Item{}
	tasksTitles := utils.GetListFromTreeTask(tasksData)
	for _, task := range tasksTitles {
		var title string = ""
		if len(task.Data.Title) > 40 {
			title = task.Data.Title[0:40]
		} else {
			title = task.Data.Title
		}
		i := item(fmt.Sprintf("%s | %s | %s ", task.Hash, task.Data.State, title))
		items = append(items, i)
	}

	l := list.New(items, itemDelegate{}, 20, 14)
	l.Title = "Jobs tasks :"

	// Initialize the viewport

	return &RootModel{
		tasks:     l,
		tasksData: tasksData,
		data:      data,
		jid:       jid,
	}
}

func (r RootModel) Init() tea.Cmd {
	return nil
}

func (r RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
		// Calculate sizes for split view (20/80)
		listWidth := r.width * 20 / 100
		viewportWidth := r.width*80/100 - 4 // subtract padding

		// Update list width
		r.tasks.SetWidth(listWidth)
		r.tasks.SetHeight(r.height - 10) // subtract space for header and section

		// Update viewport width
		r.logViewport = viewport.New(viewportWidth, r.height-10)
		r.logViewport.SetContent(r.logContent)
		//r.logViewport.Width = viewportWidth
		//r.logViewport.Height = r.height - 10
		//r.logViewport.SetHeight(r.height - 10) // subtract space for header and section
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return r, tea.Quit
		case "enter":
			i := r.tasks.Cursor()
			fmt.Println("Selected i ", i)
			foobar := r.tasksData[i]
			fmt.Println(foobar)
			taskLog := requests.GetTaskLog(r.data["user"].(string), r.jid, "9")
			r.logViewport.SetContent(taskLog)
		}
	}
	// Handle list updates
	newTasks, cmd := r.tasks.Update(msg)
	r.tasks = newTasks
	cmds = append(cmds, cmd)

	// Handle viewport updates
	newviewport, cmd := r.logViewport.Update(msg)
	r.logViewport = newviewport
	cmds = append(cmds, cmd)

	return r, tea.Batch(cmds...)
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
	// Split view (list and viewport)
	splitView := lipgloss.JoinHorizontal(
		lipgloss.Left,
		r.style.BorderStyle.Render(r.tasks.View()),  // 20%
		containerStyle.Render(r.logViewport.View()), // 80%
	)

	// Join all sections vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		jobData,
		splitView,
	)
	//return lipgloss.Place(r.width, r.height, lipgloss.Center, lipgloss.Top, lipgloss.JoinVertical(lipgloss.Top, header, jobData))
}

// func Show(data, tasksData map[string]any) {
func Show(data map[string]any, tasksData []any, jid string) {

	//main := &RootModel{data: data}
	main := initModel(data, tasksData, jid)
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
