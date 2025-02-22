package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/rebay1982/hflogger/pkg/ansi"

	tea "github.com/charmbracelet/bubbletea"
)

// KNOWN EDGE CASE
// When the log wraps around (appendIndex > bufferSize) and the window index (and window range) is smaller than the
// appendIndex - bufferSize, the user start seeing newly appended messages overwrite the oldest message. Haven't
// figured out a behaviour I'm satisfied with so this will stay this way until then (ie, prob never).
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
		case "t":
			l.tail = true
			l.tailLog()
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
	if l.selectedIndex < (l.appendIndex - 1) {
		l.selectedIndex++

		// Automatically enable taliing only if the selected line is the last line in the log.
		if l.selectedIndex == (l.appendIndex - 1) {
			l.tail = true
		}
	}
}

func (l *Log) decSelectedIndex() {
	if l.selectedIndex > 0 && (l.selectedIndex > (l.appendIndex - len(l.lines))) {
		l.selectedIndex--
	}

	// When decrementing the selected index, we'll never be at the last line so disable tailing.
	l.tail = false
}

func (l *Log) tailLog() {
	l.selectedIndex = l.appendIndex - 1
	windowEndIndex := l.windowIndex + (l.displayLines - 1)

	// if we're tailing update the window start position
	if l.selectedIndex > windowEndIndex {
		l.windowIndex++
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
