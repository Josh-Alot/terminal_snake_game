package input

import (
	"io"

	"github.com/terminal_snake_game/internal/game"
)

func StartInputLoop(r io.Reader, vimMode bool, dirCh chan<- game.Direction, quitCh chan<- struct{}) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			key, err := ReadKey(r)
			if err != nil {
				return
			}
			if key == "\x03" || key == "q" {
				quitCh <- struct{}{}
				return
			}
			dir, ok := MapKey(key, vimMode)
			if ok {
				dirCh <- dir
			}
		}
	}()
	return done
}
