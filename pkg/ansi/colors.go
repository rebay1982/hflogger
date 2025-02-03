package ansi

import (
	"fmt"
)

const (
	clear = "\033[0m"

	black   = "\033[30m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"

	boldBlack   = "\033[30;1m"
	boldRed     = "\033[31;1m"
	boldGreen   = "\033[32;1m"
	boldYellow  = "\033[33;1m"
	boldBlue    = "\033[34;1m"
	boldMagenta = "\033[35;1m"
	boldCyan    = "\033[36;1m"
	boldWhite   = "\033[37;1m"
)

func Black(str string) string {
	return fmt.Sprintf("%s%s%s", black, str, clear)
}

func Red(str string) string {
	return fmt.Sprintf("%s%s%s", red, str, clear)
}

func Green(str string) string {
	return fmt.Sprintf("%s%s%s", green, str, clear)
}

func Yellow(str string) string {
	return fmt.Sprintf("%s%s%s", yellow, str, clear)
}

func Blue(str string) string {
	return fmt.Sprintf("%s%s%s", blue, str, clear)
}

func Magenta(str string) string {
	return fmt.Sprintf("%s%s%s", magenta, str, clear)
}

func Cyan(str string) string {
	return fmt.Sprintf("%s%s%s", cyan, str, clear)
}

func White(str string) string {
	return fmt.Sprintf("%s%s%s", white, str, clear)
}

func BoldBlack(str string) string {
	return fmt.Sprintf("%s%s%s", boldBlack, str, clear)
}

func BoldRed(str string) string {
	return fmt.Sprintf("%s%s%s", boldRed, str, clear)
}

func BoldGreen(str string) string {
	return fmt.Sprintf("%s%s%s", boldGreen, str, clear)
}

func BoldYellow(str string) string {
	return fmt.Sprintf("%s%s%s", boldYellow, str, clear)
}

func BoldBlue(str string) string {
	return fmt.Sprintf("%s%s%s", boldBlue, str, clear)
}

func BoldMagenta(str string) string {
	return fmt.Sprintf("%s%s%s", boldMagenta, str, clear)
}

func BoldCyan(str string) string {
	return fmt.Sprintf("%s%s%s", boldCyan, str, clear)
}

func BoldWhite(str string) string {
	return fmt.Sprintf("%s%s%s", boldWhite, str, clear)
}
