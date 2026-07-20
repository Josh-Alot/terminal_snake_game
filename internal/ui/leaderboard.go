package ui

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type Score struct {
	Name  string
	Score int
}

func RenderLeaderboard(w io.Writer, size Size, scores []Score) {
	ClearScreen(w)

	lines := []string{"L E A D E R B O A R D", "", "R E T U R N   [ Q ]"}
	if len(scores) == 0 {
		lines = append(lines, "", "No scores yet!")
	} else {
		for i, s := range scores {
			lines = append(lines, fmt.Sprintf("%2d. %s %3d", i+1, s.Name, s.Score))
		}
	}

	width := 0
	for _, l := range lines {
		if n := utf8.RuneCountInString(l); n > width {
			width = n
		}
	}
	width += 4

	var frame strings.Builder
	topPad := (size.Height - len(lines)) / 2
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
