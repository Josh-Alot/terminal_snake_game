# Agent Instructions

## Project Goal

Create a Terminal Snake Game

## Game Requirements

### Rendering & Grid

- Dynamic grid size based on terminal dimensions.
- ASCII characters: `█` for snake body, `*` for food.
- Flicker-free rendering using ANSI escape codes (e.g., `\033[H`).

### Input Handling

- Terminal "raw mode" for instant key capture without Enter.
- Controls: `WASD`, `Arrow Keys`, and `hjkl` (Vim mode).
- Vim mode toggleable from the start screen.
- Prevent 180-degree turns (ignore opposite direction input).

### Game Logic & State

- Snake state: Slice of coordinate structs.
- Food: Randomly spawned, must not overlap with snake body.
- Difficulty: Game speed increases (tick interval decreases) as the snake grows.
- Game Over: Triggered by wall or self-collision.

### Game Flow & UI

- Start Screen (LazyVim style):
  - Centered ASCII art title.
  - Navigable menu: `Start Game`, `Leaderboard`, `Vim Mode: [ON/OFF]`, `Quit`.
- Player Identity: 3-character uppercase name required at Game Over.
- Game Over Screen: ASCII box with score, `R` to restart, `Q` to return to menu.
- Leaderboard: Persistent top 10 scores saved in `leaderboard.txt`.

## Go Style & Best Practices

- Use `gofmt` for all source files.
- Follow idiomatic Go: explicit error handling, concise naming.
- Keep game state encapsulated in a struct.
- Use a channel-based input loop to avoid blocking the game tick.
- Avoid external dependencies unless necessary for terminal manipulation.

## Testing

- Always prioritize unit tests when implementing features (TDD-first approach).
- Follow a strict per-function test => implementation cycle: write ONE test, run it and watch it fail (red), implement the minimum code to make it pass, run it again (green), then move on to the next test. Never write all tests up front and implement everything afterwards in a batch.
- Abstract system/terminal dependencies behind interfaces so they can be mocked in unit tests.
- Reserve integration/smoke tests (requiring a real TTY) for verifying actual OS-level behavior.
- Run `go test ./...` and `gofmt ./...` before considering a feature complete.
