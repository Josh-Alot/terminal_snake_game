package ui

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

// NameInput captures the 3-character uppercase player name shown on the
// game over screen.
type NameInput struct {
	chars []byte // 0..3, uppercase A-Z
}

// AddChar appends c, uppercasing lowercase letters. It returns false
// for non-letters or when the 3-character name is already complete.
func (n *NameInput) AddChar(c byte) bool {
	if len(n.chars) >= 3 {
		return false
	}
	if c >= 'a' && c <= 'z' {
		c -= 'a' - 'A'
	}
	if c < 'A' || c > 'Z' {
		return false
	}
	n.chars = append(n.chars, c)
	return true
}

// Backspace removes the last char; it is a no-op when empty.
func (n *NameInput) Backspace() {
	if len(n.chars) > 0 {
		n.chars = n.chars[:len(n.chars)-1]
	}
}

// Clear removes all entered chars.
func (n *NameInput) Clear() {
	n.chars = n.chars[:0]
}

func (n *NameInput) Value() string {
	return string(n.chars)
}

func (n *NameInput) IsComplete() bool {
	return len(n.chars) == 3
}

func (n *NameInput) Len() int {
	return len(n.chars)
}

// RenderGameOver clears the screen and draws a centered ASCII box with
// the final score and the 3-character name field (empty slots shown as
// placeholders, so the cursor position is implicit).
func RenderGameOver(w io.Writer, size Size, score int, name *NameInput, hint string) {
	ClearScreen(w)

	nameField := name.Value() + strings.Repeat("_", 3-name.Len())

	lines := []string{
		"GAME OVER",
		"",
		fmt.Sprintf("Score: %d", score),
		"",
		"Name: " + nameField,
		"",
		"[R] Restart   [Q] Menu",
	}
	if hint != "" {
		lines = append(lines, "", hint)
	}

	width := 0
	for _, l := range lines {
		if n := utf8.RuneCountInString(l); n > width {
			width = n
		}
	}
	width += 4 // 2 borders + 2 spaces of inner padding

	var frame strings.Builder

	topPad := (size.Height - (len(lines) + 2)) / 2
	if topPad < 0 {
		topPad = 0
	}
	for i := 0; i < topPad; i++ {
		frame.WriteString("\r\n")
	}

	leftPad := (size.Width - width) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	pad := strings.Repeat(" ", leftPad)

	frame.WriteString(pad + "┌" + strings.Repeat("─", width-2) + "┐\r\n")
	for _, l := range lines {
		frame.WriteString(pad + "│ " + l + strings.Repeat(" ", width-4-utf8.RuneCountInString(l)) + " │\r\n")
	}
	frame.WriteString(pad + "└" + strings.Repeat("─", width-2) + "┘\r\n")

	io.WriteString(w, frame.String())
}
