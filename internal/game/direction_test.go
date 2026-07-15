package game

import "testing"

func TestOpposite(t *testing.T) {
	tests := []struct {
		name string
		d    Direction
		want Direction
	}{
		{"Up to Down", Up, Down},
		{"Down to Up", Down, Up},
		{"Left to Right", Left, Right},
		{"Right to Left", Right, Left},
		{"Zero to Zero", Zero, Zero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.d.Opposite()
			if got != tt.want {
				t.Errorf("%s.Opposite() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestIsValidTurn(t *testing.T) {
	tests := []struct {
		name    string
		current Direction
		next    Direction
		want    bool
	}{
		{"same direction", Up, Up, true},
		{"180 degree turn up-down", Up, Down, false},
		{"180 degree turn down-up", Down, Up, false},
		{"180 degree turn left-right", Left, Right, false},
		{"180 degree turn right-left", Right, Left, false},
		{"perpendicular up-right", Up, Right, true},
		{"perpendicular up-left", Up, Left, true},
		{"perpendicular down-right", Down, Right, true},
		{"perpendicular down-left", Down, Left, true},
		{"from zero any direction", Zero, Up, true},
		{"from zero to zero", Zero, Zero, true},
		{"to zero from direction", Up, Zero, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidTurn(tt.current, tt.next)
			if got != tt.want {
				t.Errorf("IsValidTurn(%v, %v) = %v, want %v",
					tt.current, tt.next, got, tt.want)
			}
		})
	}
}
