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
	displayWidth int
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
	log.displayWidth = 80
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
	// -7 is the number of characters in the title bar that are NOT part of the title string.
	output := fmt.Sprintf("---[ %s ]%s\n", ansi.BoldGreen(l.title), strings.Repeat("-", l.displayWidth-7-len(l.title)))

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
	percentString := fmt.Sprintf("%d%%", l.computeViewportPercentage())
	output += fmt.Sprintf("%s[ %s ]---\n", strings.Repeat("-", l.displayWidth-7-len(percentString)), percentString)

	return output
}

func (l Log) computeViewportPercentage() int {
	minIndex := l.appendIndex - len(l.lines)
	if minIndex < 0 {
		minIndex = 0
	}

	maxIndex := l.appendIndex - l.displayLines
	if maxIndex < 0 {
		maxIndex = 0
	}

	// Edge case
	if minIndex == maxIndex {
		return 100
	}
	percent := float64(l.windowIndex-minIndex) / float64(maxIndex-minIndex) * 100.0

	return int(percent)
}

func (l *Log) Push(message string) {
	currentTime := time.Now().Format("2006-01-02T15:04:05.000 -0700")

	index := l.getBufferIndex()
	l.lines[index] = fmt.Sprintf("%s: %s", string(currentTime), message)

	l.appendIndex++

	if l.tail {
		l.tailLog()
	}

	// Adjust the window when the appending loops back to the beginning of the buffer.
	l.adjustWindow()
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

func (l *Log) adjustWindow() {
	bufferSize := len(l.lines)
	windowDiff := (l.appendIndex - l.windowIndex) - bufferSize

	if windowDiff > 0 {
		l.windowIndex += windowDiff
		l.selectedIndex += windowDiff
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
