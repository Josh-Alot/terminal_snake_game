package input

import (
	"io"

	"github.com/terminal_snake_game/internal/game"
)

func ReadKey(r io.Reader) (string, error) {
	buf := make([]byte, 1)
	_, err := r.Read(buf)
	if err != nil {
		return "", err
	}

	if buf[0] == 0x1b {
		esc := make([]byte, 2)
		_, err := r.Read(esc)
		if err != nil {
			return string(buf), nil
		}
		return string(append(buf, esc...)), nil
	}

	return string(buf), nil
}

func MapKey(key string, vimMode bool) (game.Direction, bool) {
	switch key {
	case "w", "W":
		return game.Up, true
	case "a", "A":
		return game.Left, true
	case "s", "S":
		return game.Down, true
	case "d", "D":
		return game.Right, true
	case "\x1b[A":
		return game.Up, true
	case "\x1b[B":
		return game.Down, true
	case "\x1b[C":
		return game.Right, true
	case "\x1b[D":
		return game.Left, true
	}

	if vimMode {
		switch key {
		case "h":
			return game.Left, true
		case "j":
			return game.Down, true
		case "k":
			return game.Up, true
		case "l":
			return game.Right, true
		}
	}

	return game.Zero, false
}
