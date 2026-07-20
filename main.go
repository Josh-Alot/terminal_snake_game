package main

import (
	"io"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/terminal_snake_game/internal/game"
	"github.com/terminal_snake_game/internal/input"
	"github.com/terminal_snake_game/internal/storage"
	"github.com/terminal_snake_game/internal/ui"
)

// Action is the result of a screen: what the main flow should do next.
type Action int

const (
	ActionStart Action = iota
	ActionQuit
	ActionRestart
	ActionMenu
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

	size, err := ui.GetTerminalSize()
	if err != nil {
		log.Fatalf("failed to get the terminal size: %v", err)
	}

	leaderboard, err := storage.NewLeaderboard("leaderboard.txt")
	if err != nil {
		log.Fatalf("failed to initialize leaderboard: %v", err)
	}

	// The terminal size is fixed for the whole session; resizing during
	// play is a known limitation (see PHASE4.md §4.7).
	gridW := (size.Width - 2) / 2
	gridH := size.Height - 2
	if gridW < 10 || gridH < 10 {
		log.Fatal("terminal too small: need at least 22x12")
	}

	// A single input goroutine feeds the whole app. Per-screen input
	// loops (input.StartInputLoop) cannot be cancelled while blocked on
	// Read, so a leftover goroutine would steal keys from the next
	// screen after a game ends.
	keyCh := make(chan string, 16)
	go readKeys(os.Stdin, keyCh)

	vimMode := false
	for {
		action := runStartScreen(size, &vimMode, keyCh, leaderboard)
		if action == ActionQuit {
			return
		}
		for {
			score, quit := runGame(size, vimMode, keyCh)
			if quit {
				return
			}
			action := runGameOver(size, score, keyCh, leaderboard)
			if action != ActionRestart {
				break
			}
		}
	}
}

// readKeys forwards raw key tokens to keyCh for the app's lifetime. On
// read error (e.g. stdin closed) it synthesises Ctrl+C so the app exits
// cleanly instead of blocking forever.
func readKeys(r io.Reader, keyCh chan<- string) {
	for {
		key, err := input.ReadKey(r)
		if err != nil {
			keyCh <- "\x03"
			return
		}
		keyCh <- key
	}
}

// drainKeys discards keys buffered during a screen transition so stale
// input is not replayed on the new screen.
func drainKeys(keyCh chan string) {
	for {
		select {
		case <-keyCh:
		default:
			return
		}
	}
}

// runStartScreen renders the menu and handles navigation until the
// player starts a game or quits. It mutates vimMode via the toggle.
func runStartScreen(size ui.Size, vimMode *bool, keyCh chan string, lb *storage.Leaderboard) Action {
	menu := ui.NewStartMenu(*vimMode)
	drainKeys(keyCh)
	ui.RenderStartScreen(os.Stdout, size, menu)

	for {
		key := <-keyCh
		switch ui.MapMenuKey(key) {
		case ui.MenuUp:
			menu.MoveUp()
		case ui.MenuDown:
			menu.MoveDown()
		case ui.MenuQuit:
			return ActionQuit
		case ui.MenuSelect:
			switch menu.Selected().ID {
			case ui.IDStart:
				*vimMode = menu.IsVimMode()
				return ActionStart
			case ui.IDQuit:
				return ActionQuit
			case ui.IDVimToggle:
				menu.ToggleVim()
			case ui.IDLeaderboard:
				drainKeys(keyCh)
				showLeaderboardScreen(size, lb, keyCh)
				drainKeys(keyCh)
			}
		}
		ui.RenderStartScreen(os.Stdout, size, menu)
	}
}

// showLeaderboardScreen renders the top scores until the user dismisses
// the screen with q. It does not return an Action: the caller stays in
// its own loop and re-renders the start screen afterwards, so the
// leaderboard can never leak into the game loop.
func showLeaderboardScreen(size ui.Size, lb *storage.Leaderboard, keyCh chan string) {
	scores := lb.TopScores()
	uiScores := make([]ui.Score, len(scores))
	for i, s := range scores {
		uiScores[i] = ui.Score{Name: s.Name, Score: s.Score}
	}

	ui.RenderLeaderboard(os.Stdout, size, uiScores)

	for {
		key := <-keyCh
		if key == "\x03" || key == "q" {
			return
		}
	}
}

// runGame plays one game and returns the final score. quit is true when
// the player asked to exit the app (q / Ctrl+C) instead of dying.
func runGame(size ui.Size, vimMode bool, keyCh chan string) (score int, quit bool) {
	// Each logical grid cell is rendered as 2 horizontal characters, so
	// gridW is half the number of playable columns.
	gridW := (size.Width - 2) / 2
	gridH := size.Height - 2
	boxSize := ui.Size{Width: 2*gridW + 2, Height: gridH + 2}

	ui.ClearScreen(os.Stdout)

	fb := ui.NewFrameBuffer(boxSize.Width, boxSize.Height)
	rng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	gs := game.NewGameState(gridW, gridH)

	ticker := time.NewTicker(time.Duration(gs.TickMs) * time.Millisecond)
	defer ticker.Stop()

	drainKeys(keyCh)
	renderFrame(os.Stdout, fb, boxSize, gs)

	for {
		select {
		case key := <-keyCh:
			if key == "\x03" || key == "q" {
				return gs.Score, true
			}
			if dir, ok := input.MapKey(key, vimMode); ok {
				gs.SetDirection(dir)
			}
		case <-ticker.C:
			gs.Update(rng)
			if gs.IsGameOver() {
				return gs.Score, false
			}
			ticker.Reset(time.Duration(gs.TickMs) * time.Millisecond)
			renderFrame(os.Stdout, fb, boxSize, gs)
		}
	}
}

// runGameOver shows the game over screen, captures the 3-letter player
// name and returns ActionRestart or ActionMenu.
func runGameOver(size ui.Size, score int, keyCh chan string, lb *storage.Leaderboard) Action {
	var name ui.NameInput
	hint := ""

	drainKeys(keyCh)
	ui.RenderGameOver(os.Stdout, size, score, &name, hint)

	for {
		key := <-keyCh
		switch key {
		case "\x03", "q":
			return ActionMenu
		case "r":
			if name.IsComplete() {
				lb.AddScore(name.Value(), score)
				if err := lb.Save(); err != nil {
					log.Printf("failed to save leaderboard: %v", err)
				}
				return ActionRestart
			}
			hint = "Enter 3 letters"
		case "\x7f":
			name.Backspace()
			hint = ""
		default:
			if len(key) == 1 && name.AddChar(key[0]) {
				hint = ""
			}
		}
		ui.RenderGameOver(os.Stdout, size, score, &name, hint)
	}
}

// renderFrame builds one frame in the framebuffer and flushes it; only
// the cells that changed since the previous tick are written. Each
// logical cell occupies a 2x1 block of terminal characters (roughly
// square in pixels), keeping the perceived speed equal on both axes.
func renderFrame(w io.Writer, fb *ui.FrameBuffer, boxSize ui.Size, gs *game.GameState) {
	fb.Clear()
	fb.DrawBox(boxSize)
	for _, p := range gs.Snake {
		fb.WriteString(2*p.X, p.Y+1, "██")
	}
	fb.WriteString(2*gs.Food.X, gs.Food.Y+1, "**")
	fb.Flush(w)
}
