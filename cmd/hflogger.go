package main

import (
	//"encoding/json"
	"fmt"

	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/rebay1982/hflogger/internal/ui"
	//"github.com/rebay1982/wsjtx-udp"
)

//var quit = make(chan struct{})
//var messageCh = make(chan wsjtxudp.WSJTXMessage)

func main() {
	app := ui.InitializeApplication("- HFLogger -")

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Couldn't run Tea program: %v\n", err)
		os.Exit(1)
	}

	//close(quit)
}
