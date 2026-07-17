package main

import (
	"io"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"time"

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

	// Each logical grid cell is rendered as 2 horizontal characters, so
	// gridW is half the number of playable columns.
	gridW := (size.Width - 2) / 2
	gridH := size.Height - 2
	if gridW < 10 || gridH < 10 {
		log.Fatal("terminal too small: need at least 22x12")
	}

	// Effective box size: 2 columns per logical cell plus borders.
	boxSize := ui.Size{Width: 2*gridW + 2, Height: gridH + 2}
	if err := ui.DrawBox(os.Stdout, boxSize); err != nil {
		log.Fatalf("failed to draw box: %v", err)
	}

	rng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	gs := game.NewGameState(gridW, gridH)

	dirCh := make(chan game.Direction, 10)
	quitCh := make(chan struct{}, 1)
	done := input.StartInputLoop(os.Stdin, false, dirCh, quitCh)

	ticker := time.NewTicker(time.Duration(gs.TickMs) * time.Millisecond)
	defer ticker.Stop()

	render(os.Stdout, gs)

gameloop:
	for {
		select {
		case dir := <-dirCh:
			gs.SetDirection(dir)
		case <-ticker.C:
			gs.Update(rng)
			if gs.IsGameOver() {
				break gameloop
			}
			ticker.Reset(time.Duration(gs.TickMs) * time.Millisecond)
			render(os.Stdout, gs)
		case <-quitCh:
			break gameloop
		case <-done:
			break gameloop
		}
	}
}

// render redraws the whole frame (box, snake and food) in a single
// buffered write to minimise flicker. Each logical cell occupies a 2x1
// block of terminal characters (roughly square in pixels), which keeps
// the perceived speed equal on both axes.
func render(w io.Writer, gs *game.GameState) {
	var buf strings.Builder

	size := ui.Size{Width: 2*gs.GridWidth + 2, Height: gs.GridHeight + 2}
	ui.DrawBox(&buf, size)

	for _, p := range gs.Snake {
		ui.DrawAt(&buf, 2*p.X, p.Y+1, "██")
	}
	ui.DrawAt(&buf, 2*gs.Food.X, gs.Food.Y+1, "**")

	w.Write([]byte(buf.String()))
}
