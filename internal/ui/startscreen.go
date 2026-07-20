package ui

import (
	"io"
	"strings"
	"unicode/utf8"
)

var titleArt = []string{
	`███████╗███╗   ██╗ █████╗ ██╗  ██╗███████╗`,
	`██╔════╝████╗  ██║██╔══██╗██║ ██╔╝██╔════╝`,
	`███████╗██╔██╗ ██║███████║█████╔╝ █████╗  `,
	`╚════██║██║╚██╗██║██╔══██║██╔═██╗ ██╔══╝  `,
	`███████║██║ ╚████║██║  ██║██║  ██╗███████╗`,
	`╚══════╝╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝`,
}

const titleMenuGap = 6

// RenderStartScreen clears the screen and draws the centered ASCII title
// above the vertically centered menu.
func RenderStartScreen(w io.Writer, size Size, menu *Menu) {
	ClearScreen(w)

	items := menu.Items()
	// blank line between each menu item → 2*n-1 content lines for n items
	menuLines := len(items)*2 - 1
	if menuLines < 0 {
		menuLines = 0
	}
	lines := make([]string, 0, len(titleArt)+titleMenuGap+menuLines)
	lines = append(lines, titleArt...)
	for i := 0; i < titleMenuGap; i++ {
		lines = append(lines, "")
	}
	for i, item := range items {
		if i > 0 {
			lines = append(lines, "")
		}
		prefix := "  "
		if i == menu.SelectedIndex() {
			prefix = "▶ "
		}
		label := spaceOut(item.Label)
		if item.ID == IDVimToggle {
			if menu.IsVimMode() {
				label += "   [ O N ]"
			} else {
				label += "   [ O F F ]"
			}
		}
		lines = append(lines, prefix+label)
	}

	topPad := (size.Height - len(lines)) / 2
	if topPad < 0 {
		topPad = 0
	}

	var frame strings.Builder
	for i := 0; i < topPad; i++ {
		frame.WriteString("\r\n")
	}
	for _, line := range lines {
		writeCentered(&frame, size.Width, line)
	}

	io.WriteString(w, frame.String())
}

// spaceOut uppercases s, inserts a space between each letter, and
// triple-spaces between words.
func spaceOut(s string) string {
	words := strings.Fields(strings.ToUpper(s))
	spaced := make([]string, len(words))
	for i, w := range words {
		runes := []rune(w)
		parts := make([]string, len(runes))
		for j, r := range runes {
			parts[j] = string(r)
		}
		spaced[i] = strings.Join(parts, " ")
	}
	return strings.Join(spaced, "   ")
}

func writeCentered(frame *strings.Builder, width int, line string) {
	pad := (width - utf8.RuneCountInString(line)) / 2
	if pad < 0 {
		pad = 0
	}
	frame.WriteString(strings.Repeat(" ", pad))
	frame.WriteString(line)
	frame.WriteString("\r\n")
}
