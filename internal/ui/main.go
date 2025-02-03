package ui

import (
	"fmt"

	"github.com/rebay1982/hflogger/internal/server"
	"github.com/rebay1982/hflogger/pkg/ansi"

	tea "github.com/charmbracelet/bubbletea"
)

type MainApplication struct {
	title     string
	cmdBar    tea.Model
	logWindow tea.Model
	udpServer *server.WSJTXServer
}

func InitializeApplication(title string) MainApplication {
	logWindow, _ := NewLog("WSJT-X LOG", 5, 1000)
	server, _ := server.NewWSJTXServer("127.0.0.1", 2237)

	return MainApplication{
		title:     title,
		cmdBar:    NewCommandBar(),
		logWindow: logWindow,
		udpServer: server,
	}
}

func (m MainApplication) Init() tea.Cmd {

	return m.getMsgFromServer
}

func (m MainApplication) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.logWindow, _ = m.logWindow.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.udpServer.Close()
			return m, tea.Quit
		}
	}

	return m, m.getMsgFromServer
}

func (m MainApplication) View() string {
	output := fmt.Sprintf("%s\n\n", ansi.BoldWhite(m.title))

	output += m.logWindow.View()
	output += m.cmdBar.View()

	return output
}

type logMessage string
type errMessage string

func (m MainApplication) getMsgFromServer() tea.Msg {
	message, err := m.udpServer.ReadFromUDP()
	if err != nil {
		return errMessage(err.Error())
	}

	return logMessage(message.Header.MsgType.String())
}
