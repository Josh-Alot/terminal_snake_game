package game

import "testing"

func TestPoint_Equality(t *testing.T) {
	tests := []struct {
		name string
		a    Point
		b    Point
		want bool
	}{
		{"same X and Y", Point{X: 3, Y: 4}, Point{X: 3, Y: 4}, true},
		{"different X", Point{X: 3, Y: 4}, Point{X: 4, Y: 4}, false},
		{"different Y", Point{X: 3, Y: 4}, Point{X: 3, Y: 5}, false},
		{"different X and Y", Point{X: 3, Y: 4}, Point{X: 5, Y: 6}, false},
		{"zero values", Point{}, Point{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a == tt.b
			if got != tt.want {
				t.Errorf("Point%v == Point%v = %v, want %v",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}
