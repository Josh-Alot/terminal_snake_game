package game

import (
	"math/rand/v2"
	"testing"
)

func newTestRNG() *rand.Rand {
	return rand.New(rand.NewPCG(1, 2))
}

func newTestState() *GameState {
	gs := NewGameState(40, 20)
	// Move food away from the snake's path so movement tests are
	// not affected by random food placement.
	gs.Food = Point{X: 1, Y: 1}
	return gs
}

func snakeContains(snake []Point, p Point) bool {
	for _, s := range snake {
		if s == p {
			return true
		}
	}
	return false
}

func TestNewGameState(t *testing.T) {
	gs := NewGameState(40, 20)

	if len(gs.Snake) != 3 {
		t.Fatalf("snake length = %d, want 3", len(gs.Snake))
	}

	head := Point{X: 40/2 + 1, Y: 20 / 2}
	if gs.Snake[0] != head {
		t.Errorf("head = %v, want %v", gs.Snake[0], head)
	}
	if gs.Snake[1] != (Point{X: 20, Y: 10}) {
		t.Errorf("body = %v, want {20 10}", gs.Snake[1])
	}
	if gs.Snake[2] != (Point{X: 19, Y: 10}) {
		t.Errorf("tail = %v, want {19 10}", gs.Snake[2])
	}

	if gs.Direction != Right {
		t.Errorf("direction = %v, want Right", gs.Direction)
	}
	if gs.Score != 0 {
		t.Errorf("score = %d, want 0", gs.Score)
	}
	if gs.TickMs != 150 {
		t.Errorf("tickMs = %d, want 150", gs.TickMs)
	}
	if gs.IsGameOver() {
		t.Error("IsGameOver() = true, want false")
	}
	if gs.Food == (Point{}) {
		t.Error("food was not placed (zero position)")
	}
	if gs.Food.X < 1 || gs.Food.X > 40 || gs.Food.Y < 1 || gs.Food.Y > 20 {
		t.Errorf("food %v out of grid bounds", gs.Food)
	}
	if snakeContains(gs.Snake, gs.Food) {
		t.Errorf("food %v overlaps snake body", gs.Food)
	}
}

func TestUpdate_MovesRight(t *testing.T) {
	gs := newTestState()
	gs.Update(newTestRNG())

	if gs.Snake[0] != (Point{X: 22, Y: 10}) {
		t.Errorf("head = %v, want {22 10}", gs.Snake[0])
	}
	if len(gs.Snake) != 3 {
		t.Errorf("snake length = %d, want 3", len(gs.Snake))
	}
	if gs.Snake[2] != (Point{X: 20, Y: 10}) {
		t.Errorf("tail = %v, want {20 10}", gs.Snake[2])
	}
}

func TestUpdate_MovesUp(t *testing.T) {
	gs := newTestState()
	gs.SetDirection(Up)
	gs.Update(newTestRNG())

	if gs.Direction != Up {
		t.Errorf("direction = %v, want Up", gs.Direction)
	}
	if gs.Snake[0] != (Point{X: 21, Y: 9}) {
		t.Errorf("head = %v, want {21 9}", gs.Snake[0])
	}
	if len(gs.Snake) != 3 {
		t.Errorf("snake length = %d, want 3", len(gs.Snake))
	}
}

func TestUpdate_MovesDown(t *testing.T) {
	gs := newTestState()
	gs.SetDirection(Down)
	gs.Update(newTestRNG())

	if gs.Direction != Down {
		t.Errorf("direction = %v, want Down", gs.Direction)
	}
	if gs.Snake[0] != (Point{X: 21, Y: 11}) {
		t.Errorf("head = %v, want {21 11}", gs.Snake[0])
	}
	if len(gs.Snake) != 3 {
		t.Errorf("snake length = %d, want 3", len(gs.Snake))
	}
}

func TestUpdate_MovesLeft(t *testing.T) {
	gs := newTestState()

	// Left is a 180-degree turn from Right, so turn Up first and
	// then Left on the following tick.
	gs.SetDirection(Up)
	gs.Update(newTestRNG())
	gs.SetDirection(Left)
	gs.Update(newTestRNG())

	if gs.Direction != Left {
		t.Errorf("direction = %v, want Left", gs.Direction)
	}
	if gs.Snake[0] != (Point{X: 20, Y: 9}) {
		t.Errorf("head = %v, want {20 9}", gs.Snake[0])
	}
	if len(gs.Snake) != 3 {
		t.Errorf("snake length = %d, want 3", len(gs.Snake))
	}
}

func TestUpdate_GrowsOnFood(t *testing.T) {
	gs := newTestState()
	gs.Food = Point{X: 22, Y: 10} // directly in front of the head

	gs.Update(newTestRNG())

	if len(gs.Snake) != 4 {
		t.Errorf("snake length = %d, want 4", len(gs.Snake))
	}
	if gs.Snake[0] != (Point{X: 22, Y: 10}) {
		t.Errorf("head = %v, want {22 10}", gs.Snake[0])
	}
	if gs.Score != 1 {
		t.Errorf("score = %d, want 1", gs.Score)
	}
	if !gs.grew {
		t.Error("grew = false, want true")
	}
	// A new food must have been spawned on an empty cell.
	if gs.Food == (Point{X: 22, Y: 10}) {
		t.Error("food was not respawned after being eaten")
	}
	if snakeContains(gs.Snake, gs.Food) {
		t.Errorf("new food %v overlaps snake body", gs.Food)
	}
}

func TestUpdate_NoGrowWithoutFood(t *testing.T) {
	gs := newTestState()
	gs.Update(newTestRNG())

	if len(gs.Snake) != 3 {
		t.Errorf("snake length = %d, want 3", len(gs.Snake))
	}
	if gs.Score != 0 {
		t.Errorf("score = %d, want 0", gs.Score)
	}
	if gs.grew {
		t.Error("grew = true, want false")
	}
}

func TestUpdate_WallCollision(t *testing.T) {
	tests := []struct {
		name      string
		snake     []Point
		direction Direction
	}{
		{"right wall", []Point{{10, 5}, {9, 5}, {8, 5}}, Right},
		{"left wall", []Point{{1, 5}, {2, 5}, {3, 5}}, Left},
		{"top wall", []Point{{5, 1}, {5, 2}, {5, 3}}, Up},
		{"bottom wall", []Point{{5, 10}, {5, 9}, {5, 8}}, Down},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GameState{
				Snake:      tt.snake,
				Direction:  tt.direction,
				Food:       Point{X: 10, Y: 10},
				TickMs:     150,
				GridWidth:  10,
				GridHeight: 10,
			}

			gs.Update(newTestRNG())

			if !gs.IsGameOver() {
				t.Error("IsGameOver() = false, want true")
			}
			if gs.Snake[0] != tt.snake[0] {
				t.Errorf("head moved to %v after wall collision, want %v",
					gs.Snake[0], tt.snake[0])
			}
		})
	}
}

func TestUpdate_SelfCollision(t *testing.T) {
	// Snake coiled so that moving Down drives the head into its own body.
	gs := &GameState{
		Snake: []Point{
			{2, 2}, {2, 1}, {1, 1}, {1, 2}, {1, 3}, {2, 3}, {3, 3},
		},
		Direction:  Down,
		Food:       Point{X: 10, Y: 10},
		TickMs:     150,
		GridWidth:  10,
		GridHeight: 10,
	}

	gs.Update(newTestRNG())

	if !gs.IsGameOver() {
		t.Error("IsGameOver() = false, want true")
	}
}

func TestUpdate_SelfCollision_TailRemoved(t *testing.T) {
	// Moving Down places the head exactly on the tail's current cell,
	// which is vacated on the same tick: no collision.
	gs := &GameState{
		Snake: []Point{
			{3, 2}, {3, 1}, {2, 1}, {1, 1}, {1, 2}, {1, 3}, {2, 3}, {3, 3},
		},
		Direction:  Down,
		Food:       Point{X: 10, Y: 10},
		TickMs:     150,
		GridWidth:  10,
		GridHeight: 10,
	}

	gs.Update(newTestRNG())

	if gs.IsGameOver() {
		t.Error("IsGameOver() = true, want false (snake can follow its own tail)")
	}
	if gs.Snake[0] != (Point{X: 3, Y: 3}) {
		t.Errorf("head = %v, want {3 3}", gs.Snake[0])
	}
	if len(gs.Snake) != 8 {
		t.Errorf("snake length = %d, want 8", len(gs.Snake))
	}
}

func TestSetDirection_StoresPending(t *testing.T) {
	gs := newTestState()
	gs.SetDirection(Up)

	if gs.pendingDir != Up {
		t.Errorf("pendingDir = %v, want Up", gs.pendingDir)
	}
}

func TestSetDirection_OverwritesOnMultipleCalls(t *testing.T) {
	gs := newTestState()
	gs.SetDirection(Up)
	gs.SetDirection(Down)

	if gs.pendingDir != Down {
		t.Errorf("pendingDir = %v, want Down (last call wins)", gs.pendingDir)
	}
}

func TestSetDirection_Rejected180_InUpdate(t *testing.T) {
	gs := newTestState()
	gs.SetDirection(Left) // opposite of Right
	gs.Update(newTestRNG())

	if gs.Direction != Right {
		t.Errorf("direction = %v, want Right (180-degree turn rejected)", gs.Direction)
	}
	if gs.pendingDir != Zero {
		t.Errorf("pendingDir = %v, want Zero after Update", gs.pendingDir)
	}
	if gs.Snake[0] != (Point{X: 22, Y: 10}) {
		t.Errorf("head = %v, want {22 10} (kept moving Right)", gs.Snake[0])
	}
}

func TestSpawnFood_PlacesOnEmpty(t *testing.T) {
	gs := NewGameState(40, 20)

	for i := 0; i < 200; i++ {
		rng := rand.New(rand.NewPCG(uint64(i), 0))
		gs.SpawnFood(rng)

		if gs.Food.X < 1 || gs.Food.X > 40 || gs.Food.Y < 1 || gs.Food.Y > 20 {
			t.Fatalf("iteration %d: food %v out of grid bounds", i, gs.Food)
		}
		if snakeContains(gs.Snake, gs.Food) {
			t.Fatalf("iteration %d: food %v overlaps snake body", i, gs.Food)
		}
	}
}

func TestSpawnFood_FullBoardGameOver(t *testing.T) {
	gs := &GameState{GridWidth: 4, GridHeight: 4}
	for y := 1; y <= 4; y++ {
		for x := 1; x <= 4; x++ {
			gs.Snake = append(gs.Snake, Point{X: x, Y: y})
		}
	}

	gs.SpawnFood(newTestRNG())

	if !gs.IsGameOver() {
		t.Error("IsGameOver() = false, want true (board full)")
	}
}

func TestSpawnFood_DeterministicWithRNG(t *testing.T) {
	gs1 := NewGameState(40, 20)
	gs2 := NewGameState(40, 20)

	rng1 := rand.New(rand.NewPCG(42, 7))
	rng2 := rand.New(rand.NewPCG(42, 7))

	gs1.SpawnFood(rng1)
	gs2.SpawnFood(rng2)

	if gs1.Food != gs2.Food {
		t.Errorf("same seed produced different food positions: %v vs %v",
			gs1.Food, gs2.Food)
	}
}

func TestTickInterval_DecreasesWithScore(t *testing.T) {
	gs := NewGameState(200, 50)

	for score := 1; score <= 20; score++ {
		head := gs.Snake[0]
		gs.Food = Point{X: head.X + 1, Y: head.Y}
		gs.Update(newTestRNG())

		want := 150 - score*5
		if want < 50 {
			want = 50
		}
		if gs.TickMs != want {
			t.Errorf("score %d: TickMs = %d, want %d", score, gs.TickMs, want)
		}
		if gs.Score != score {
			t.Errorf("Score = %d, want %d", gs.Score, score)
		}
	}
}

func TestTickInterval_CappedAt50(t *testing.T) {
	gs := NewGameState(200, 50)

	for i := 1; i <= 25; i++ {
		head := gs.Snake[0]
		gs.Food = Point{X: head.X + 1, Y: head.Y}
		gs.Update(newTestRNG())

		if gs.TickMs < 50 {
			t.Fatalf("score %d: TickMs = %d, must never go below 50", i, gs.TickMs)
		}
	}

	if gs.Score != 25 {
		t.Errorf("Score = %d, want 25", gs.Score)
	}
	if gs.TickMs != 50 {
		t.Errorf("TickMs = %d, want 50 (capped)", gs.TickMs)
	}
}
