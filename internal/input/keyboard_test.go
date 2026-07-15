package input

import (
	"strings"
	"testing"

	"github.com/terminal_snake_game/internal/game"
)

func TestMapKey_WASD(t *testing.T) {
	tests := []struct {
		key     string
		want    game.Direction
		wantOk  bool
		vimMode bool
	}{
		{"w", game.Up, true, false},
		{"W", game.Up, true, false},
		{"a", game.Left, true, false},
		{"A", game.Left, true, false},
		{"s", game.Down, true, false},
		{"S", game.Down, true, false},
		{"d", game.Right, true, false},
		{"D", game.Right, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, ok := MapKey(tt.key, tt.vimMode)
			if got != tt.want || ok != tt.wantOk {
				t.Errorf("MapKey(%q, %v) = (%v, %v), want (%v, %v)",
					tt.key, tt.vimMode, got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func TestMapKey_ArrowKeys(t *testing.T) {
	tests := []struct {
		key  string
		want game.Direction
	}{
		{"\x1b[A", game.Up},
		{"\x1b[B", game.Down},
		{"\x1b[C", game.Right},
		{"\x1b[D", game.Left},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, ok := MapKey(tt.key, false)
			if !ok || got != tt.want {
				t.Errorf("MapKey(%q, false) = (%v, %v), want (%v, true)",
					tt.key, got, ok, tt.want)
			}
		})
	}
}

func TestMapKey_VimMode(t *testing.T) {
	tests := []struct {
		key     string
		want    game.Direction
		wantOk  bool
		vimMode bool
	}{
		{"h", game.Left, true, true},
		{"j", game.Down, true, true},
		{"k", game.Up, true, true},
		{"l", game.Right, true, true},
		{"h", game.Zero, false, false},
		{"j", game.Zero, false, false},
		{"k", game.Zero, false, false},
		{"l", game.Zero, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, ok := MapKey(tt.key, tt.vimMode)
			if got != tt.want || ok != tt.wantOk {
				t.Errorf("MapKey(%q, %v) = (%v, %v), want (%v, %v)",
					tt.key, tt.vimMode, got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func TestMapKey_InvalidKey(t *testing.T) {
	tests := []string{"x", " ", "1", "\x1b[Z", ""}

	for _, key := range tests {
		t.Run(key, func(t *testing.T) {
			got, ok := MapKey(key, false)
			if ok || got != game.Zero {
				t.Errorf("MapKey(%q, false) = (%v, %v), want (Zero, false)",
					key, got, ok)
			}
		})
	}
}

func TestReadKey_SingleByte(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"w key", "w", "w", false},
		{"W key", "W", "W", false},
		{"a key", "a", "a", false},
		{"space", " ", " ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadKey(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadKey() err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ReadKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadKey_ArrowKeys(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"up arrow", "\x1b[A", "\x1b[A"},
		{"down arrow", "\x1b[B", "\x1b[B"},
		{"right arrow", "\x1b[C", "\x1b[C"},
		{"left arrow", "\x1b[D", "\x1b[D"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadKey(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("ReadKey() unexpected err: %v", err)
			}
			if got != tt.want {
				t.Errorf("ReadKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReadKey_EmptyReader(t *testing.T) {
	_, err := ReadKey(strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error on empty reader, got nil")
	}
}
