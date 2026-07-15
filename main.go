package main

import (
	"fmt"
	"log"
	"os"

	"github.com/terminal_snake_game/internal/game"
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
	defer func() {
		ui.ClearScreen(os.Stdout)
		ui.MoveCursor(os.Stdout, 1, 1)
		ui.ShowCursor(os.Stdout)
	}()

	ui.ClearScreen(os.Stdout)
	ui.MoveCursor(os.Stdout, 1, 1)

	size, err := ui.GetTerminalSize()
	if err != nil {
		log.Fatalf("failed to get the terminal size: %v", err)
	}

	if err := ui.DrawBox(os.Stdout, size); err != nil {
		log.Fatalf("failed to draw box: %v", err)
	}

	dirCh := make(chan game.Direction, 10)
	quitCh := make(chan struct{}, 1)
	done := input.StartInputLoop(os.Stdin, false, dirCh, quitCh)

	for {
		select {
		case dir := <-dirCh:
			ui.MoveCursor(os.Stdout, 2, 2)
			fmt.Fprintf(os.Stdout, "direction: %v   ", dir)
		case <-quitCh:
			return
		case <-done:
			return
		}
	}
}
