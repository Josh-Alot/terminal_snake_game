package ui

import "testing"

func TestMapMenuKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want MenuEvent
	}{
		{"arrow up", "\x1b[A", MenuUp},
		{"k moves up", "k", MenuUp},
		{"arrow down", "\x1b[B", MenuDown},
		{"j moves down", "j", MenuDown},
		{"enter selects", "\r", MenuSelect},
		{"space selects", " ", MenuSelect},
		{"q quits", "q", MenuQuit},
		{"ctrl+c quits", "\x03", MenuQuit},
		{"w is unmapped", "w", MenuNone},
		{"s is unmapped", "s", MenuNone},
		{"arrow left is unmapped", "\x1b[D", MenuNone},
		{"unknown key", "x", MenuNone},
		{"empty key", "", MenuNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MapMenuKey(tt.key); got != tt.want {
				t.Errorf("MapMenuKey(%q) = %d, want %d", tt.key, got, tt.want)
			}
		})
	}
}
