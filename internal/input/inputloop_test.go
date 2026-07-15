package input

import (
	"strings"
	"testing"
	"time"

	"github.com/terminal_snake_game/internal/game"
)

func TestStartInputLoop_SendsMappedDirections(t *testing.T) {
	keys := "wsd"
	r := strings.NewReader(keys)
	dirCh := make(chan game.Direction, 10)
	quitCh := make(chan struct{}, 1)

	done := StartInputLoop(r, false, dirCh, quitCh)
	defer close(dirCh)

	got := collectDirections(dirCh, 3, time.Second)

	if len(got) != 3 {
		t.Fatalf("expected 3 directions, got %d", len(got))
	}
	if got[0] != game.Up {
		t.Errorf("expected Up, got %v", got[0])
	}
	if got[1] != game.Down {
		t.Errorf("expected Down, got %v", got[1])
	}
	if got[2] != game.Right {
		t.Errorf("expected Right, got %v", got[2])
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("input loop did not finish after reader exhausted")
	}
}

func TestStartInputLoop_VimMode(t *testing.T) {
	keys := "kj"
	r := strings.NewReader(keys)
	dirCh := make(chan game.Direction, 10)
	quitCh := make(chan struct{}, 1)

	StartInputLoop(r, true, dirCh, quitCh)
	defer close(dirCh)

	got := collectDirections(dirCh, 2, time.Second)

	if len(got) != 2 {
		t.Fatalf("expected 2 directions, got %d", len(got))
	}
	if got[0] != game.Up {
		t.Errorf("expected Up, got %v", got[0])
	}
	if got[1] != game.Down {
		t.Errorf("expected Down, got %v", got[1])
	}
}

func TestStartInputLoop_IgnoresUnmappedKeys(t *testing.T) {
	keys := "wx"
	r := strings.NewReader(keys)
	dirCh := make(chan game.Direction, 10)
	quitCh := make(chan struct{}, 1)

	StartInputLoop(r, false, dirCh, quitCh)
	defer close(dirCh)

	got := collectDirections(dirCh, 1, time.Second)

	if len(got) != 1 {
		t.Fatalf("expected 1 direction (ignored unmapped), got %d", len(got))
	}
	if got[0] != game.Up {
		t.Errorf("expected Up, got %v", got[0])
	}
}

func TestStartInputLoop_QuitOnCtrlC(t *testing.T) {
	r := strings.NewReader("\x03")
	dirCh := make(chan game.Direction, 10)
	quitCh := make(chan struct{}, 1)

	StartInputLoop(r, false, dirCh, quitCh)
	defer close(dirCh)

	select {
	case <-quitCh:
	case <-time.After(time.Second):
		t.Fatal("expected quit signal on Ctrl+C, timed out")
	}
}

func TestStartInputLoop_QuitOnQ(t *testing.T) {
	r := strings.NewReader("q")
	dirCh := make(chan game.Direction, 10)
	quitCh := make(chan struct{}, 1)

	StartInputLoop(r, false, dirCh, quitCh)
	defer close(dirCh)

	select {
	case <-quitCh:
	case <-time.After(time.Second):
		t.Fatal("expected quit signal on 'q', timed out")
	}
}

func collectDirections(ch <-chan game.Direction, count int, timeout time.Duration) []game.Direction {
	var result []game.Direction
	deadline := time.After(timeout)
	for len(result) < count {
		select {
		case d := <-ch:
			result = append(result, d)
		case <-deadline:
			return result
		}
	}
	return result
}
