package ui

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

func DrawAt(writer io.Writer, x, y int, ch string) {
	MoveCursor(writer, x, y)
	fmt.Fprint(writer, ch)
}

func DrawBox(writer io.Writer, size Size) error {
	if size.Width < 3 || size.Height < 3 {
		return errors.New("terminal size too small to draw a box")
	}

	var frame strings.Builder

	top := "┌" + strings.Repeat("─", size.Width-2) + "┐"
	mid := "│" + strings.Repeat(" ", size.Width-2) + "│"
	bot := "└" + strings.Repeat("─", size.Width-2) + "┘"

	DrawAt(&frame, 1, 1, top)
	for y := 2; y < size.Height; y++ {
		DrawAt(&frame, 1, y, mid)
	}

	DrawAt(&frame, 1, size.Height, bot)

	_, err := writer.Write([]byte(frame.String()))
	return err
}
