# Implementation Roadmap: Terminal Snake Game

## Phase 1: Foundation & Terminal Control

- [X] **Environment Setup**: Initialize Go module and establish a clean project structure (`/internal/game`, `/internal/ui`, `/internal/input`).
- [X] **Terminal Raw Mode**: Implement a mechanism to put the terminal into "raw mode" to capture keystrokes instantly without requiring the Enter key.
- [X] **ANSI Utility Library**: Build a set of helper functions for:
  - `MoveCursor(x, y)`: Using `\033[H` and related codes.
  - `ClearScreen()`: Efficiently clearing the view.
  - `HideCursor()` / `ShowCursor()`: Improving visual polish.
- [X] **Dimension Tracking**: Implement dynamic detection of terminal width and height to define the game grid boundaries.
- [X] **Basic Rendering**: Build a simple `Draw` primitive to write characters at specific coordinates, validating the ANSI utilities and dimension tracking (e.g., drawing a static bordered box).

## Phase 2: The Input Engine

- [ ] **Concurrent Input Loop**: Launch a dedicated goroutine to listen for keyboard events.
- [ ] **Input Mapping**:
  - Map `WASD` and `Arrow Keys` to directions.
  - Implement `hjkl` mapping for Vim mode.
- [ ] **Input Channel**: Use a thread-safe channel to send movement commands to the main game loop, ensuring the game tick doesn't block on user input.
- [ ] **Directional Logic**: Implement a guard to prevent "180-degree turns" (e.g., if moving North, ignore South input).

## Phase 3: Core Game Logic

- [ ] **State Modeling**:
  - `Point` struct for `X, Y` coordinates.
  - `GameState` struct to track the snake (slice of `Point`), food position, current direction, score, and tick interval.
- [ ] **The Game Loop**: Implement a ticker-based loop:
    1. **Process Input**: Check the input channel for direction changes.
    2. **Update State**: Move the snake head, shift the body, and check for food.
    3. **Collision Detection**: Check if the head intersects with walls or the snake's own body.
    4. **Spawn Food**: Randomly generate food coordinates, ensuring they do not overlap with the snake's current body.
- [ ] **Difficulty Scaling**: Implement a formula to decrease the tick interval (increase speed) as the snake's length increases.
- [ ] **Game State Rendering**: Draw the current `GameState` (snake body and food) each tick using the Phase 1 primitives, enabling visual verification of the game logic.

## Phase 4: UI & User Experience

- [ ] **Start Screen**:
  - Render a centered ASCII art title.
  - Create a navigable menu for `Start Game`, `Leaderboard`, `Vim Mode Toggle`, and `Quit`.
- [ ] **Rendering Pipeline**: Optimize the Phase 3 rendering into a flicker-free pipeline by calculating the full frame and writing it to the terminal in one go or using strategic cursor movement.
- [ ] **Game Over Sequence**:
  - Display an ASCII box with the final score.
  - **Player Identity**: Implement a 3-character uppercase name input field.
  - Menu options to `Restart (R)` or `Return to Menu (Q)`.

## Phase 5: Persistence & Polishing

- [ ] **Leaderboard System**:
  - Implement file I/O for `leaderboard.txt`.
  - Logic to insert new scores, sort them, and truncate the list to the top 10.
- [ ] **Visual Refinement**: Fine-tune the ASCII characters (`█` for body, `*` for food) and ensure the layout is centered regardless of terminal size.
- [ ] **Code Quality**:
  - Run `gofmt` across the project.
  - Ensure strict error handling for all terminal and file operations.
- [ ] **Verification**: Systematic testing of all three input schemes and boundary collision edge cases.
