package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rebay1982/hflogger/pkg/ansi"

	tea "github.com/charmbracelet/bubbletea"
)

type Log struct {
	displayLines int
	appendIndex  int
	lines        []string
	title        string

	selectedIndex int
	windowIndex   int
	tail          bool
}

func NewLog(title string, displayLines, bufferSize int) (Log, error) {
	log := Log{}

	if displayLines > bufferSize {
		return log, fmt.Errorf("Cannot initialize Log UI component: display lines [%d] larger than buffer size [%d]",
			displayLines,
			bufferSize)
	}

	log.displayLines = displayLines
	log.appendIndex = 0
	log.lines = make([]string, bufferSize)
	log.title = title

	log.selectedIndex = 0
	log.windowIndex = 0
	log.tail = true

	if log.title == "" {
		log.title = "LOG"
	}

	// TODO: Remove this, it's just some test data.
	for i := 0; i < 50; i++ {
		log.Push("Line number: " + strconv.Itoa(i))
	}

	return log, nil
}

func (l Log) Init() tea.Cmd {
	return nil
}

func (l Log) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "k", "up":
			l.decSelectedIndex()
		case "j", "down":
			l.incSelectedIndex()
		}
	case logMessage:
		l.Push(string(msg))
	}

	l.updateWindowIndex()

	return l, nil
}

func (l Log) View() string {
	// Header
	title := l.title
	output := fmt.Sprintf("---[ %s ]%s\n", ansi.BoldGreen(title), strings.Repeat("-", 73-len(title)))

	linesToDisplay := l.displayLines
	if l.appendIndex < linesToDisplay {
		linesToDisplay = l.appendIndex
	}

	for i := 0; i < linesToDisplay; i++ {
		index := l.getOffsetWindowStartIndex(i)

		logLine := l.lines[index]
		logCursor := " "
		if index == l.getSelectedIndex() {
			logLine = ansi.BoldWhite(logLine)
			logCursor = ansi.BoldWhite(">")
		}
		output += fmt.Sprintf("%s %s\n", logCursor, logLine)

	}

	for i := linesToDisplay; i < l.displayLines; i++ {
		output += fmt.Sprintln("")
	}

	// Footer
	output += strings.Repeat("-", 80)
	output += "\n"

	return output
}

func (l Log) getSelectedIndex() int {
	return l.selectedIndex % len(l.lines)
}

func (l Log) getWindowStartIndex() int {
	return l.getOffsetWindowStartIndex(0)
}

func (l Log) getOffsetWindowStartIndex(offset int) int {
	index := (l.windowIndex + offset) % len(l.lines)

	// TODO: I think there's a bug here. It will always repeat the last line of the buffer.
	if index < 0 {
		index = len(l.lines) - 1
	}

	return index
}

func (l Log) getBufferIndex() int {
	return l.getOffsetBufferIndex(0)
}

func (l Log) getOffsetBufferIndex(offset int) int {
	index := (l.appendIndex + offset) % len(l.lines)

	// TODO: I think there's a bug here. It will always repeat the last line of the buffer.
	if index < 0 {
		index = len(l.lines) - 1
	}

	return index
}

func (l *Log) Push(message string) {
	currentTime := time.Now().Format("2006-01-02T15:04:05.000 -0700")

	index := l.getBufferIndex()
	l.lines[index] = fmt.Sprintf("%s: %s", string(currentTime), message)

	if l.tail {
		l.selectedIndex = l.appendIndex
		windowEndIndex := l.windowIndex + (l.displayLines -1)

		// if we're tailing update the window start position
		if l.selectedIndex > windowEndIndex {
			l.windowIndex++
		}
	}

	l.appendIndex++
}

func (l *Log) incSelectedIndex() {
	l.selectedIndex++

	if l.selectedIndex == l.appendIndex {
		l.selectedIndex--
	}

}

func (l *Log) decSelectedIndex() {
	l.selectedIndex--

	if l.selectedIndex < 0 {
		l.selectedIndex++
	}
}

func (l *Log) updateWindowIndex() {
	if l.selectedIndex < l.windowIndex {
		l.windowIndex = l.selectedIndex
	}

	if l.selectedIndex > (l.windowIndex + (l.displayLines - 1)) {
		l.windowIndex = l.selectedIndex - (l.displayLines - 1)
	}
}
