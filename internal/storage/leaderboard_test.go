package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestNewLeaderboard_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "leaderboard.txt")

	lb, err := NewLeaderboard(file)
	if err != nil {
		t.Fatalf("NewLeaderboard() error = %v", err)
	}

	if len(lb.scores) != 0 {
		t.Errorf("expected empty leaderboard, got %d entries", len(lb.scores))
	}
}

func TestLoadLeaderboard_ParsesValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "leaderboard.txt")

	content := "ABC 100\r\nDEF 95\r\nGHI 90\r\n"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	lb, err := NewLeaderboard(file)
	if err != nil {
		t.Fatalf("NewLeaderboard() error = %v", err)
	}

	if len(lb.scores) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(lb.scores))
	}

	if lb.scores[0].Name != "ABC" || lb.scores[0].Score != 100 {
		t.Errorf("first entry = {%s, %d}, want {ABC, 100}", lb.scores[0].Name, lb.scores[0].Score)
	}
}

func TestLoadLeaderboard_CreatesFileIfMissing(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "leaderboard.txt")

	_, err := NewLeaderboard(file)
	if err != nil {
		t.Fatalf("NewLeaderboard() error = %v", err)
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Error("expected leaderboard file to be created")
	}
}

func TestAddScore_InsertsSorted(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "leaderboard.txt")

	lb, err := NewLeaderboard(file)
	if err != nil {
		t.Fatalf("NewLeaderboard() error = %v", err)
	}

	lb.AddScore("ABC", 90)
	lb.AddScore("DEF", 100)
	lb.AddScore("GHI", 85)

	if len(lb.scores) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(lb.scores))
	}

	if lb.scores[0].Name != "DEF" || lb.scores[0].Score != 100 {
		t.Errorf("first entry = {%s, %d}, want {DEF, 100}", lb.scores[0].Name, lb.scores[0].Score)
	}
	if lb.scores[1].Name != "ABC" || lb.scores[1].Score != 90 {
		t.Errorf("second entry = {%s, %d}, want {ABC, 90}", lb.scores[1].Name, lb.scores[1].Score)
	}
	if lb.scores[2].Name != "GHI" || lb.scores[2].Score != 85 {
		t.Errorf("third entry = {%s, %d}, want {GHI, 85}", lb.scores[2].Name, lb.scores[2].Score)
	}
}

func TestAddScore_TruncatesToTen(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "leaderboard.txt")

	lb, err := NewLeaderboard(file)
	if err != nil {
		t.Fatalf("NewLeaderboard() error = %v", err)
	}

	for i := 0; i < 15; i++ {
		lb.AddScore(fmt.Sprintf("PLY%02d", i), 100-i)
	}

	if len(lb.scores) != 10 {
		t.Errorf("expected 10 entries (truncated), got %d", len(lb.scores))
	}

	if lb.scores[9].Score != 91 {
		t.Errorf("10th score = %d, want 91", lb.scores[9].Score)
	}
}

func TestAddScore_IgnoresLowerScores(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "leaderboard.txt")

	lb, err := NewLeaderboard(file)
	if err != nil {
		t.Fatalf("NewLeaderboard() error = %v", err)
	}

	lb.AddScore("ABC", 100)
	lb.AddScore("ABC", 50)

	if len(lb.scores) != 1 {
		t.Errorf("expected 1 entry (lower score ignored), got %d", len(lb.scores))
	}
}

func TestSave_WritesCorrectFormat(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "leaderboard.txt")

	lb, err := NewLeaderboard(file)
	if err != nil {
		t.Fatalf("NewLeaderboard() error = %v", err)
	}

	lb.AddScore("ABC", 100)
	lb.AddScore("DEF", 95)

	if err := lb.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "ABC 100\r\nDEF 95\r\n"
	if string(content) != expected {
		t.Errorf("file content = %q, want %q", string(content), expected)
	}
}

func TestTopScores_ReturnsCopy(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "leaderboard.txt")

	lb, err := NewLeaderboard(file)
	if err != nil {
		t.Fatalf("NewLeaderboard() error = %v", err)
	}

	lb.AddScore("ABC", 100)
	lb.AddScore("DEF", 95)

	scores := lb.TopScores()

	if len(scores) != 2 {
		t.Fatalf("expected 2 scores, got %d", len(scores))
	}

	scores[0].Score = 200
	if lb.scores[0].Score == 200 {
		t.Error("TopScores() should return a copy, not a reference")
	}
}
