package game

type Direction int

const (
	Zero Direction = iota
	Up
	Down
	Left
	Right
)

func (d Direction) Opposite() Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	default:
		return Zero
	}
}

func IsValidTurn(current, next Direction) bool {
	if next == Zero {
		return true
	}
	return next != current.Opposite()
}
