package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestDrawAt_WritesCharAtPosition(t *testing.T) {
	var buffer bytes.Buffer

	DrawAt(&buffer, 5, 3, "X")

	want := "\033[3;5H"
	got := buffer.String()

	if !strings.HasPrefix(got, want) {
		t.Errorf("expected %q, got %q", want, got)
	}

	if !strings.HasSuffix(got, "X") {
		t.Errorf("expected X, got %q", got)
	}
}

func TestDrawBox_DrawsCornersAndBorders(t *testing.T) {
	var buffer bytes.Buffer

	if err := DrawBox(&buffer, Size{Width: 5, Height: 3}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buffer.String()
	lines := strings.Split(got, "\n")

	if !strings.Contains(lines[0], "┌") || !strings.Contains(lines[0], "┐") {
		t.Errorf("expected a top border with corners, got %q", lines[0])
	}

	if !strings.Contains(got, "└") || !strings.Contains(got, "┘") {
		t.Errorf("expected bottom border with corners, got %q", got)
	}

	if strings.Count(got, "│") < 2 {
		t.Errorf("expected at least 2 side borders '│', got %q", got)
	}
}

func TestDrawBox_RejectsTooSmallSize(t *testing.T) {
	var buffer bytes.Buffer

	err := DrawBox(&buffer, Size{Width: 2, Height: 2})
	if err == nil {
		t.Fatal("expected an error when size is too small for a box, got a box")
	}
}
