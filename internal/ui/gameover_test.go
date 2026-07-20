package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestNameInput_AddLetter(t *testing.T) {
	var n NameInput

	if !n.AddChar('a') {
		t.Fatal("expected AddChar to accept a letter")
	}
	if !n.AddChar('B') {
		t.Fatal("expected AddChar to accept an uppercase letter")
	}

	if got := n.Value(); got != "AB" {
		t.Errorf("expected lowercase input to be uppercased, got %q", got)
	}
}

func TestNameInput_RejectsNonLetter(t *testing.T) {
	var n NameInput

	for _, c := range []byte{'1', '9', '!', ' ', '-', '\r'} {
		if n.AddChar(c) {
			t.Errorf("expected AddChar(%q) to be rejected", c)
		}
	}

	if got := n.Value(); got != "" {
		t.Errorf("expected no chars to be appended, got %q", got)
	}
}

func TestNameInput_CapsAtThree(t *testing.T) {
	var n NameInput

	n.AddChar('A')
	n.AddChar('B')
	n.AddChar('C')

	if n.AddChar('D') {
		t.Error("expected AddChar to return false when full")
	}
	if got := n.Value(); got != "ABC" {
		t.Errorf("expected value to stay %q, got %q", "ABC", got)
	}
}

func TestNameInput_Backspace(t *testing.T) {
	var n NameInput

	n.AddChar('A')
	n.AddChar('B')
	n.Backspace()

	if got := n.Value(); got != "A" {
		t.Errorf("expected %q after backspace, got %q", "A", got)
	}

	n.Backspace()
	n.Backspace() // no-op at length 0

	if got := n.Value(); got != "" {
		t.Errorf("expected empty value, got %q", got)
	}
}

func TestNameInput_IsComplete(t *testing.T) {
	var n NameInput

	if n.IsComplete() {
		t.Error("expected empty input to be incomplete")
	}

	for _, c := range []byte{'A', 'B', 'C'} {
		n.AddChar(c)
		if n.Len() > 3 {
			t.Fatalf("expected Len() <= 3, got %d", n.Len())
		}
	}

	if n.Len() != 3 {
		t.Errorf("expected Len() 3, got %d", n.Len())
	}
	if !n.IsComplete() {
		t.Error("expected input with 3 chars to be complete")
	}
}

func TestNameInput_Clear(t *testing.T) {
	var n NameInput

	n.AddChar('A')
	n.AddChar('B')
	n.Clear()

	if n.Len() != 0 {
		t.Errorf("expected Len() 0 after Clear, got %d", n.Len())
	}
	if got := n.Value(); got != "" {
		t.Errorf("expected empty value after Clear, got %q", got)
	}
}

func TestRenderGameOver_ContainsScore(t *testing.T) {
	var buf bytes.Buffer
	var name NameInput

	RenderGameOver(&buf, Size{Width: 80, Height: 24}, 42, &name, "")

	if got := buf.String(); !strings.Contains(got, "Score: 42") {
		t.Errorf("expected output to contain %q, got:\n%s", "Score: 42", got)
	}
}

func TestRenderGameOver_NameField(t *testing.T) {
	tests := []struct {
		name    string
		letters string
		want    string
	}{
		{"no letters shows all placeholders", "", "___"},
		{"two letters pads the rest", "AB", "AB_"},
		{"full name has no placeholders", "ABC", "ABC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var name NameInput
			for i := 0; i < len(tt.letters); i++ {
				name.AddChar(tt.letters[i])
			}

			RenderGameOver(&buf, Size{Width: 80, Height: 24}, 0, &name, "")

			if got := buf.String(); !strings.Contains(got, tt.want) {
				t.Errorf("expected name field to contain %q, got:\n%s", tt.want, got)
			}
		})
	}
}

func TestRenderGameOver_Hints(t *testing.T) {
	var buf bytes.Buffer
	var name NameInput

	RenderGameOver(&buf, Size{Width: 80, Height: 24}, 7, &name, "Enter 3 letters")

	got := buf.String()
	for _, want := range []string{"[R] Restart", "[Q] Menu", "Enter 3 letters"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, got)
		}
	}
}

func TestRenderGameOver_UsesCRLFForRawMode(t *testing.T) {
	var buf bytes.Buffer
	var name NameInput

	RenderGameOver(&buf, Size{Width: 80, Height: 24}, 42, &name, "")

	got := buf.String()
	if !strings.Contains(got, "GAME OVER\r\n") && !strings.Contains(got, "GAME OVER") {
		t.Fatalf("expected GAME OVER in output, got:\n%q", got)
	}
	// Content after clear must use \r\n so lines stay left-aligned in raw mode.
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
