package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderStartScreen_ContainsTitle(t *testing.T) {
	var buf bytes.Buffer
	menu := NewStartMenu(false)

	RenderStartScreen(&buf, Size{Width: 80, Height: 30}, menu)

	if got := buf.String(); !strings.Contains(got, `███████╗███╗   ██╗ █████╗ ██╗  ██╗███████╗`) {
		t.Errorf("expected output to contain the block ASCII title, got:\n%s", got)
	}
}

func TestRenderStartScreen_ContainsAllItems(t *testing.T) {
	var buf bytes.Buffer
	menu := NewStartMenu(false)

	RenderStartScreen(&buf, Size{Width: 80, Height: 30}, menu)

	got := buf.String()
	for _, label := range []string{
		"S T A R T   G A M E",
		"L E A D E R B O A R D",
		"V I M   M O D E",
		"Q U I T",
	} {
		if !strings.Contains(got, label) {
			t.Errorf("expected output to contain item %q, got:\n%s", label, got)
		}
	}
}

func TestRenderStartScreen_VimModeState(t *testing.T) {
	tests := []struct {
		name    string
		vimMode bool
		want    string
		notWant string
	}{
		{"vim on shows [ O N ]", true, "[ O N ]", "[ O F F ]"},
		{"vim off shows [ O F F ]", false, "[ O F F ]", "[ O N ]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			RenderStartScreen(&buf, Size{Width: 80, Height: 30}, NewStartMenu(tt.vimMode))

			got := buf.String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("expected output to contain %q, got:\n%s", tt.want, got)
			}
			if strings.Contains(got, tt.notWant) {
				t.Errorf("expected output not to contain %q, got:\n%s", tt.notWant, got)
			}
		})
	}
}

func TestRenderStartScreen_UsesCRLFForRawMode(t *testing.T) {
	var buf bytes.Buffer
	menu := NewStartMenu(false)

	RenderStartScreen(&buf, Size{Width: 80, Height: 30}, menu)

	got := buf.String()
	// After clear, content lines must end with \r\n so the cursor returns to
	// column 0 in raw mode (bare \n only moves down and skews the layout).
	if !strings.Contains(got, titleArt[0]+"\r\n") {
		t.Errorf("expected title lines terminated with \\r\\n, got:\n%q", got)
	}
	if !strings.Contains(got, "▶ S T A R T   G A M E\r\n") {
		t.Errorf("expected menu lines terminated with \\r\\n, got:\n%q", got)
	}
	// No bare LF without preceding CR in the frame body (after clear codes).
	body := got
	if i := strings.Index(got, titleArt[0]); i >= 0 {
		body = got[i:]
	}
	for i := 0; i < len(body); i++ {
		if body[i] == '\n' && (i == 0 || body[i-1] != '\r') {
			t.Fatalf("found bare \\n at body offset %d (raw mode needs \\r\\n)", i)
		}
	}
}

func TestRenderStartScreen_HighlightsSelected(t *testing.T) {
	tests := []struct {
		name     string
		moveDown int
		wantLine string
	}{
		{"start selected by default", 0, "▶ S T A R T   G A M E"},
		{"vim mode selected after two moves down", 2, "▶ V I M   M O D E"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			menu := NewStartMenu(false)
			for i := 0; i < tt.moveDown; i++ {
				menu.MoveDown()
			}

			RenderStartScreen(&buf, Size{Width: 80, Height: 30}, menu)

			got := buf.String()
			if count := strings.Count(got, "▶"); count != 1 {
				t.Errorf("expected exactly one ▶, got %d in:\n%s", count, got)
			}
			if !strings.Contains(got, tt.wantLine) {
				t.Errorf("expected selected line %q, got:\n%s", tt.wantLine, got)
			}
		})
	}
}

func TestRenderStartScreen_TitleMenuGap(t *testing.T) {
	var buf bytes.Buffer
	RenderStartScreen(&buf, Size{Width: 80, Height: 30}, NewStartMenu(false))

	got := buf.String()
	after := got[strings.Index(got, titleArt[len(titleArt)-1]):]
	rest := after[len(titleArt[len(titleArt)-1]):]
	// Drop the title line's own trailing \r\n.
	if !strings.HasPrefix(rest, "\r\n") {
		t.Fatalf("expected title line to end with \\r\\n, got %q", rest[:min(20, len(rest))])
	}
	rest = rest[2:]
	// writeCentered("") emits pad spaces then \r\n
	blanks := 0
	for {
		i := 0
		for i < len(rest) && rest[i] == ' ' {
			i++
		}
		if i < len(rest) && strings.HasPrefix(rest[i:], "\r\n") {
			blanks++
			rest = rest[i+2:]
			continue
		}
		break
	}
	if blanks != titleMenuGap {
		t.Errorf("expected %d blank lines between title and menu, got %d", titleMenuGap, blanks)
	}
	if !strings.HasPrefix(strings.TrimLeft(rest, " "), "▶ S T A R T   G A M E") {
		t.Errorf("expected menu to follow the gap, got prefix %q", rest[:min(40, len(rest))])
	}
}

func TestSpaceOut(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"Start Game", "S T A R T   G A M E"},
		{"Leaderboard", "L E A D E R B O A R D"},
		{"Vim Mode", "V I M   M O D E"},
		{"Quit", "Q U I T"},
	}
	for _, tt := range tests {
		if got := spaceOut(tt.in); got != tt.want {
			t.Errorf("spaceOut(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
