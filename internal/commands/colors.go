package commands

import "fmt"

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

func colorize(color, text string) string {
	return color + text + colorReset
}

func printError(msg string) {
	fmt.Println(colorize(colorRed, "✗ "+msg))
}

func printSuccess(msg string) {
	fmt.Println(colorize(colorGreen, "✓ "+msg))
}

func printInfo(msg string) {
	fmt.Println(colorize(colorCyan, msg))
}
