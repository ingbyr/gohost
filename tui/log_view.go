package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"time"
)

func NewLogView(model *Model) *LogView {
	m := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.Title = "Logs"
	m.SetShowStatusBar(true)
	m.SetShowHelp(false)
	m.SetShowTitle(true)
	m.SetShowPagination(false)
	return &LogView{
		main:  model,
		model: m,
	}
}

type LogView struct {
	main  *Model
	model list.Model
}

func (l *LogView) Init() tea.Cmd {
	return nil
}

func (l *LogView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m := msg.(type) {
	case tea.WindowSizeMsg:
		l.model.SetWidth(m.Width)
		l.model.SetHeight(m.Height)
	}
	l.model, cmd = l.model.Update(msg)
	return l, cmd
}

func (l *LogView) View() string {
	return l.model.View()
}

func (l *LogView) InsertLog(msg string) {
	last := len(l.model.Items())
	l.model.InsertItem(last, &LogItem{
		log:  msg,
		time: time.Now().String(),
	})
	l.model.Select(last)
}

type LogItem struct {
	log  string
	time string
}

func (l LogItem) FilterValue() string {
	return l.log
}

func (l LogItem) Title() string {
	return l.log
}

func (l LogItem) Description() string {
	return l.time
}

type LogItemDelegate struct {
}

func (d *LogItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var str string
	if m.Index() == index {
		str = "> "
	} else {
		str = "  "
	}
	str += item.(*LogItem).log
	_, _ = fmt.Fprint(w, str)
}

func (d *LogItemDelegate) Height() int {
	return 1
}

func (d *LogItemDelegate) Spacing() int {
	return 1
}

func (d *LogItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
