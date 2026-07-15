package main

import (
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

	// ... rest of the game
}
