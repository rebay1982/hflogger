package ui

import (
	"github.com/rebay1982/hflogger/pkg/ansi"

	tea "github.com/charmbracelet/bubbletea"
)

type CommandBar struct {
}

func NewCommandBar() tea.Model {

	return CommandBar{}
}

func (c CommandBar) Init() tea.Cmd {
	return nil
}

func (c CommandBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return c, nil
}

func (c CommandBar) View() string {
	return ansi.Black(" q:exit\tt:tail")
}
