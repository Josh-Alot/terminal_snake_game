package game

import (
	"math/rand/v2"
	"time"
)

const (
	initialTickMs = 150
	minTickMs     = 50
	tickDecrease  = 5
)

type GameState struct {
	Snake      []Point
	Food       Point
	Direction  Direction
	Score      int
	TickMs     int
	GridWidth  int
	GridHeight int
	gameOver   bool
	pendingDir Direction // direction queued for the next tick
	grew       bool      // true when the snake just ate food
}

// NewGameState creates a game on a gridWidth x gridHeight grid. The snake
// starts with 3 segments, horizontally centered, facing Right, and the
// first food is spawned immediately.
func NewGameState(gridWidth, gridHeight int) *GameState {
	headX := gridWidth/2 + 1
	headY := gridHeight / 2

	gs := &GameState{
		Snake: []Point{
			{X: headX, Y: headY},
			{X: headX - 1, Y: headY},
			{X: headX - 2, Y: headY},
		},
		Direction:  Right,
		Score:      0,
		TickMs:     initialTickMs,
		GridWidth:  gridWidth,
		GridHeight: gridHeight,
	}

	rng := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 0))
	gs.SpawnFood(rng)

	return gs
}

func (gs *GameState) IsGameOver() bool {
	return gs.gameOver
}

// SetDirection queues d to be applied on the next tick. Turn validation
// happens in Update, so the last call before a tick always wins.
func (gs *GameState) SetDirection(d Direction) {
	gs.pendingDir = d
}

// Update advances the game by one tick.
func (gs *GameState) Update(rng *rand.Rand) {
	if gs.pendingDir != Zero && IsValidTurn(gs.Direction, gs.pendingDir) {
		gs.Direction = gs.pendingDir
	}
	gs.pendingDir = Zero

	var dx, dy int
	switch gs.Direction {
	case Up:
		dy = -1
	case Down:
		dy = 1
	case Left:
		dx = -1
	case Right:
		dx = 1
	}

	head := gs.Snake[0]
	newHead := Point{X: head.X + dx, Y: head.Y + dy}

	if newHead.X < 1 || newHead.X > gs.GridWidth ||
		newHead.Y < 1 || newHead.Y > gs.GridHeight {
		gs.gameOver = true
		return
	}

	// The tail cell is vacated this tick unless the snake just grew, so
	// the new head may safely move into it.
	body := gs.Snake
	if !gs.grew {
		body = body[:len(body)-1]
	}
	for _, p := range body {
		if p == newHead {
			gs.gameOver = true
			return
		}
	}

	gs.Snake = append([]Point{newHead}, gs.Snake...)

	if newHead == gs.Food {
		gs.grew = true
		gs.Score++
		gs.TickMs = tickInterval(gs.Score)
		gs.SpawnFood(rng)
	} else {
		gs.grew = false
		gs.Snake = gs.Snake[:len(gs.Snake)-1]
	}
}

// SpawnFood places the food on a random grid cell not occupied by the
// snake. If the board is full, the game ends (victory condition).
func (gs *GameState) SpawnFood(rng *rand.Rand) {
	occupied := make(map[Point]bool, len(gs.Snake))
	for _, p := range gs.Snake {
		occupied[p] = true
	}

	var empty []Point
	for y := 1; y <= gs.GridHeight; y++ {
		for x := 1; x <= gs.GridWidth; x++ {
			p := Point{X: x, Y: y}
			if !occupied[p] {
				empty = append(empty, p)
			}
		}
	}

	if len(empty) == 0 {
		gs.gameOver = true
		return
	}

	gs.Food = empty[rng.IntN(len(empty))]
}

// tickInterval returns the tick duration in milliseconds for a given
// score: the game speeds up as the snake grows, capped at minTickMs.
func tickInterval(score int) int {
	t := initialTickMs - score*tickDecrease
	if t < minTickMs {
		t = minTickMs
	}
	return t
}
