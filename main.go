package main

import (
	"fmt"
	"log"
	"os"

	"github.com/terminal_snake_game/internal/input"
	"github.com/terminal_snake_game/internal/ui"
)

func main() {
	restore, err := input.EnableRawMode()
	if err != nil {
		log.Fatalf("failed to enable raw mode: %v", err)
	}
	defer restore()

	ui.HideCursor(os.Stdout)
	defer ui.ShowCursor(os.Stdout)

	ui.ClearScreen(os.Stdout)
	ui.MoveCursor(os.Stdout, 1, 1)

	size, err := ui.GetTerminalSize()
	if err != nil {
		log.Fatalf("failed to get the terminal size: %v", err)
	}

	fmt.Fprintf(os.Stdout, "terminal size: %dx%d\r\n", size.Width, size.Height)
	// ... rest of the game
}
