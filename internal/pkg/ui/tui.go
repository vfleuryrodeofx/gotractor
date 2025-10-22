// // TODO :
// -
package ui

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/vfleuryrodeofx/gotractor/internal/pkg/requests"
	"github.com/vfleuryrodeofx/gotractor/internal/pkg/utils"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Application state type
type AppState int

// Enum for managing app state
const (
	tasksView AppState = iota
	logView
)

// STYLES TO FIX
// Define some basic styling
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(1, 0).
			Width(100).
			Align(lipgloss.Center) //.
		//Border(lipgloss.RoundedBorder())

	//sectionStyle      = lipgloss.NewStyle().Padding(1, 2).Width(100)
	containerStyle    = lipgloss.NewStyle().Padding(0, 0)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(0)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

// LIST Widget
type taskItem struct {
	tid         string
	jid         string
	title       string
	status      string
	description string
}

func (ti taskItem) Tid() string         { return ti.tid }
func (ti taskItem) Jid() string         { return ti.jid }
func (ti taskItem) Title() string       { return ti.title }
func (ti taskItem) Status() string      { return ti.status }
func (ti taskItem) FilterValue() string { return "" }
func (ti taskItem) Description() string {
	return fmt.Sprintf("%s | Status : %s", ti.tid, ti.status)
}

// Styles
type Styles struct {
	BorderColor         lipgloss.Color
	BorderStyle         lipgloss.Style
	BoldText            lipgloss.Style
	Underlined          lipgloss.Style
	SelectedView        lipgloss.Style
	ZoomedStyle         lipgloss.Style
	SelectedBorderColor lipgloss.Color
}

func DefaultStyle(width int) *Styles {
	s := new(Styles)
	s.SelectedBorderColor = lipgloss.Color("#e28743")
	s.BorderColor = lipgloss.Color("white")
	s.BorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(s.BorderColor).Padding(0, 1)
	s.BoldText = lipgloss.NewStyle().Bold(true)
	s.Underlined = lipgloss.NewStyle().Underline(true).Bold(true)
	s.SelectedView = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(s.SelectedBorderColor).Padding(0, 1)
	s.ZoomedStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true, false, true, false).BorderForeground(s.SelectedBorderColor).Padding(0, 1)

	return s
}

// Wrap text for viewport
func WrapText(content string, width int) string {
	style := lipgloss.NewStyle().Width(width)
	return style.Render(content)
}

// Root Model
type RootModel struct {
	state           AppState
	tasks           list.Model
	logViewport     viewport.Model
	logContent      string
	width           int
	height          int
	style           *Styles
	data            map[string]any
	tasksData       []any
	jid             string
	showOnlyLogView bool
	currentLog      string
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
		i := taskItem{tid: task.Hash, status: task.Data.State, title: title, jid: jid}
		items = append(items, i)
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
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

	//var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
		// Calculate sizes for split view (20/80)
		listWidth := r.width * 20 / 100
		//viewportWidth := r.width*80/100 - 4 // subtract padding
		viewportWidth := r.width*80/100 - 1 // subtract padding
		// Update list width
		r.tasks.SetSize(listWidth, r.height-20)
		// Update viewport width
		r.logViewport = viewport.New(viewportWidth, r.height-20)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return r, tea.Quit
		case "enter":
			task_item := r.tasks.SelectedItem()
			selectedTid := task_item.(taskItem).Tid()
			//taskLog := requests.GetTaskLog(r.data["user"].(string), r.jid, selectedTid)
			r.currentLog = requests.GetTaskLog(r.data["user"].(string), r.jid, selectedTid)
			//r.logViewport.SetContent(WrapText(taskLog, r.logViewport.Width))
			r.logViewport.SetContent(WrapText(r.currentLog, r.logViewport.Width))
		case "tab", "shift+tab":
			if r.state == tasksView {
				r.state = logView
			} else {
				r.state = tasksView
			}
		case "z":
			if r.showOnlyLogView == false {
				r.showOnlyLogView = true
				fmt.Println("Zoom !")
				r.state = logView
			} else {
				r.showOnlyLogView = false
				fmt.Println("No zoom")
			}
		}
	}

	if r.state == tasksView {
		// Handle list updates
		newTasks, cmd := r.tasks.Update(msg)
		r.tasks = newTasks
		cmds = append(cmds, cmd)
	}

	if r.state == logView {
		// Handle viewport updates
		newviewport, cmd := r.logViewport.Update(msg)
		r.logViewport = newviewport
		cmds = append(cmds, cmd)
	}
	return r, tea.Batch(cmds...)
}

func (r RootModel) View() string {
	title, ok := r.data["title"].(string)
	comment, _ := r.data["comment"].(string)
	if !ok {
		title = "No job title ..."
	}

	if r.showOnlyLogView == true {
		r.logViewport.Width = r.width
		r.logViewport.Height = r.height - 2
		r.logViewport.SetContent(WrapText(r.currentLog, r.width-1))
		zoomedView := r.style.ZoomedStyle.Render(r.logViewport.View())
		return zoomedView
	}

	header := r.style.
		BorderStyle.Width(r.width-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(r.style.BoldText.Render("🚜 GO Tractor ! 🚜\n"))
	jobData := r.style.BorderStyle.Width(r.width-2).
		Align(lipgloss.Left, lipgloss.Center).
		Render(
			r.style.Underlined.Render("\nJob title :"),
			title,
			r.style.Underlined.Render("\nProjects:"),
			r.data["projects"].([]any)[0].(string),
			r.style.Underlined.Render("\nComment :"),
			comment,
			r.style.Underlined.Render("\nUser :"),
			r.data["user"].(string),
			r.style.Underlined.Render("\nService Keys :"),
			r.data["service"].(string),
			r.style.Underlined.Render("\nSpool date :"),
			r.data["spooldate"].(string),
			r.style.Underlined.Render("\nEnv keys:"),
			r.data["envkey"].([]any)[0].(string),
		)
	var splitView string
	// Split view (list and viewport)
	if r.state == tasksView {
		splitView = lipgloss.JoinHorizontal(
			lipgloss.Top,
			r.style.SelectedView.Render(r.tasks.View()),                             // 20%
			r.style.BorderStyle.Render(containerStyle.Render(r.logViewport.View())), // 80%
		)
	} else {
		splitView = lipgloss.JoinHorizontal(
			lipgloss.Top,
			r.style.BorderStyle.Render(r.tasks.View()),                               // 20%
			r.style.SelectedView.Render(containerStyle.Render(r.logViewport.View())), // 80%
		)
	}

	// Join all sections vertically
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		jobData,
		splitView,
	)
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
	p := tea.NewProgram(*main, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
