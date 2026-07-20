package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderLeaderboard_Empty(t *testing.T) {
	var buf bytes.Buffer

	RenderLeaderboard(&buf, Size{Width: 80, Height: 24}, nil)

	got := buf.String()
	if !strings.Contains(got, "L E A D E R B O A R D") {
		t.Errorf("expected output to contain title, got:\n%s", got)
	}
	if !strings.Contains(got, "No scores yet!") {
		t.Errorf("expected output to contain 'No scores yet!', got:\n%s", got)
	}
}

func TestRenderLeaderboard_WithScores(t *testing.T) {
	var buf bytes.Buffer
	scores := []Score{
		{Name: "ABC", Score: 100},
		{Name: "DEF", Score: 95},
		{Name: "GHI", Score: 90},
	}

	RenderLeaderboard(&buf, Size{Width: 80, Height: 24}, scores)

	got := buf.String()
	for _, want := range []string{"L E A D E R B O A R D", "1. ABC 100", "2. DEF  95", "3. GHI  90", "R E T U R N"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, got)
		}
	}
}

func TestRenderLeaderboard_UsesCRLFForRawMode(t *testing.T) {
	var buf bytes.Buffer
	RenderLeaderboard(&buf, Size{Width: 80, Height: 24}, []Score{{Name: "ABC", Score: 100}})

	got := buf.String()
	body := got
	if i := strings.Index(got, "┌"); i >= 0 {
		body = got[i:]
	}
	for i := 0; i < len(body); i++ {
		if body[i] == '\n' && (i == 0 || body[i-1] != '\r') {
			t.Fatalf("found bare \\n at body offset %d (raw mode needs \\r\\n)", i)
		}
	}
}
