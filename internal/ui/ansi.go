package ui

import (
	"fmt"
	"io"
)

func MoveCursor(writer io.Writer, x, y int) {
	fmt.Fprintf(writer, "\033[%d;%dH", y, x)
}

func ClearScreen(writer io.Writer) {
	fmt.Fprintf(writer, "\033[2J\033[H")
}

func HideCursor(writer io.Writer) {
	fmt.Fprintf(writer, "\033[?25l")
}

func ShowCursor(writer io.Writer) {
	fmt.Fprint(writer, "\033[?25h")
}
