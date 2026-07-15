package ui

import (
	"bytes"
	"testing"
)

func TestMoveCursor(t *testing.T) {
	var buffer bytes.Buffer
	MoveCursor(&buffer, 3, 5)

	want := "\033[5;3H"
	got := buffer.String()

	if got != want {
		t.Errorf("moveCursor: want %s, got %s", want, got)
	}
}

func TestClearScreen(t *testing.T) {
	var buffer bytes.Buffer
	ClearScreen(&buffer)

	want := "\033[2J\033[H"
	got := buffer.String()

	if got != want {
		t.Errorf("ClearScreen: want %s, got %s", want, got)
	}
}

func TestHideCursor(t *testing.T) {
	var buffer bytes.Buffer
	HideCursor(&buffer)

	want := "\033[?25l"
	got := buffer.String()

	if got != want {
		t.Errorf("HideCursor: want %s, got %s", want, got)
	}
}

func TestShowCursor(t *testing.T) {
	var buffer bytes.Buffer
	ShowCursor(&buffer)

	want := "\033[?25h"
	got := buffer.String()

	if got != want {
		t.Errorf("ShowCursor: want %s, got %s", want, got)
	}
}
