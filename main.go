package main

import (
	"log"

	"github.com/terminal_snake_game/internal/input"
)

func main() {
	restore, err := input.EnableRawMode()
	if err != nil {
		log.Fatalf("failed to enable raw mode: %v", err)
	}
	defer restore()

	// ... rest of the game
}
