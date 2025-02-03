package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/rebay1982/hflogger/pkg/ansi"

	tea "github.com/charmbracelet/bubbletea"
)

type Log struct {
	displayLines int
	index        int
	lines        []string
	title        string
}

func NewLog(title string, displayLines, bufferSize int) (Log, error) {
	log := Log{}

	if displayLines > bufferSize {
		return log, fmt.Errorf("Cannot initialize Log UI component: display lines [%d] larger than buffer size [%d]",
			displayLines,
			bufferSize)
	}

	log.displayLines = displayLines
	log.index = 0
	log.lines = make([]string, bufferSize)
	log.title = title

	if log.title == "" {
		log.title = "LOG"
	}

	return log, nil
}

func (l Log) Init() tea.Cmd {
	return nil
}

func (l Log) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case logMessage:
		l.Push(string(msg))
	}

	return l, nil
}

func (l Log) View() string {
	// Header
	title := l.title
	output := fmt.Sprintf("---[ %s ]%s\n", ansi.BoldGreen(title), strings.Repeat("-", 73-len(title)))

	linesToDisplay := l.displayLines
	if l.index < linesToDisplay {
		linesToDisplay = l.index
	}
	for i := linesToDisplay; i > 0; i-- {
		index := l.getOffsetBufferIndex(-i)
		output += fmt.Sprintf(" %s\n", l.lines[index])
	}

	for i := linesToDisplay; i < l.displayLines; i++ {
		output += fmt.Sprintln("")
	}

	// Footer
	output += strings.Repeat("-", 80)
	output += "\n"

	return output
}

func (l *Log) Push(message string) {
	currentTime := time.Now().Format("2006-01-02T15:04:05.000 -0700")

	index := l.getBufferIndex()
	l.lines[index] = fmt.Sprintf("%s: %s", string(currentTime), message)
	l.index++
}

func (l Log) getBufferIndex() int {
	return l.getOffsetBufferIndex(0)
}

func (l Log) getOffsetBufferIndex(offset int) int {
	index := (l.index + offset) % len(l.lines)
	if index < 0 {
		index = len(l.lines) - 1
	}

	return index
}
