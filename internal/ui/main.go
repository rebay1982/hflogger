package ui

import (
	"fmt"
	"time"

	"github.com/rebay1982/hflogger/internal/server"
	"github.com/rebay1982/hflogger/pkg/ansi"

	"github.com/rebay1982/wsjtx-udp"

	tea "github.com/charmbracelet/bubbletea"
)

type MainApplication struct {
	title     string
	cmdBar    tea.Model
	logWindow tea.Model
	udpServer *server.WSJTXServer
	counter   int
}

func InitializeApplication(title string) MainApplication {
	logWindow, _ := NewLog("WSJT-X LOG", 10, 20)
	server, _ := server.NewWSJTXServer("127.0.0.1", 2237)

	return MainApplication{
		title:     title,
		cmdBar:    NewCommandBar(),
		logWindow: logWindow,
		udpServer: server,
	}
}

func (m MainApplication) Init() tea.Cmd {
	return m.getMock
	//return m.getMsgFromServer
}

func (m MainApplication) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.logWindow, _ = m.logWindow.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.udpServer.Close()
			return m, tea.Quit
		}
	case logMessage:
		return m, m.getMock
		//return m, m.getMsgFromServer
	}

	return m, nil
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

	switch message.Header.MsgType.String() {
	case "Decode":
		return logMessage(message.Payload.(wsjtxudp.DecodePayload).Message)
	default:
		return logMessage(message.Header.MsgType.String())
	}
}

var counter = 0

func (m MainApplication) getMock() tea.Msg {

	time.Sleep(1 * time.Second)
	msg := fmt.Sprintf("Hello world %d", counter)
	counter++

	return logMessage(msg)
}
