package storage

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Score struct {
	Name  string
	Score int
}

type Leaderboard struct {
	file   string
	scores []Score
}

func NewLeaderboard(file string) (*Leaderboard, error) {
	lb := &Leaderboard{file: file}

	content, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.WriteFile(file, []byte(""), 0644); err != nil {
				return nil, fmt.Errorf("failed to create leaderboard file: %w", err)
			}
			return lb, nil
		}
		return nil, fmt.Errorf("failed to read leaderboard file: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		var score int
		if _, err := fmt.Sscanf(line, "%s %d", &parts[0], &score); err != nil {
			continue
		}
		lb.scores = append(lb.scores, Score{Name: parts[0], Score: score})
	}

	return lb, nil
}

func (lb *Leaderboard) AddScore(name string, score int) {
	for i, s := range lb.scores {
		if s.Name == name && s.Score >= score {
			return
		}
		if s.Name == name {
			lb.scores[i] = Score{Name: name, Score: score}
			sort.Slice(lb.scores, func(i, j int) bool {
				return lb.scores[i].Score > lb.scores[j].Score
			})
			if len(lb.scores) > 10 {
				lb.scores = lb.scores[:10]
			}
			return
		}
	}

	lb.scores = append(lb.scores, Score{Name: name, Score: score})

	sort.Slice(lb.scores, func(i, j int) bool {
		return lb.scores[i].Score > lb.scores[j].Score
	})

	if len(lb.scores) > 10 {
		lb.scores = lb.scores[:10]
	}
}

func (lb *Leaderboard) TopScores() []Score {
	result := make([]Score, len(lb.scores))
	copy(result, lb.scores)
	return result
}

func (lb *Leaderboard) Save() error {
	var sb strings.Builder
	for i, s := range lb.scores {
		sb.WriteString(fmt.Sprintf("%s %d", s.Name, s.Score))
		if i < len(lb.scores)-1 {
			sb.WriteString("\r\n")
		}
	}

	content := sb.String()
	if len(lb.scores) > 0 {
		content += "\r\n"
	}

	if err := os.WriteFile(lb.file, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write leaderboard file: %w", err)
	}

	return nil
}
