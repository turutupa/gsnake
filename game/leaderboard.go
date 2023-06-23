package gsnake

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const FILE_NAME = "scoreboard"
const FOLDER_NAME = "gsnake"
const MAX_FILE_LINES = 100

type Score struct {
	player string
	score  int
}

type Leaderboard struct {
	scoreboardFile   string
	existsStorageDir bool
}

func NewLeaderboard() (*Leaderboard, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	// Create the directory if it doesn't exist
	scoreboardDir := filepath.Join(configDir, FOLDER_NAME)
	if _, err := os.Stat(scoreboardDir); os.IsNotExist(err) {
		err := os.MkdirAll(scoreboardDir, 0700)
		if err != nil {
			return nil, err
		}
	}

	// Create the file if it doesn't exist
	scoreboardFile := filepath.Join(scoreboardDir, FILE_NAME)
	if _, err := os.Stat(scoreboardFile); os.IsNotExist(err) {
		_, err := os.Create(scoreboardFile)
		if err != nil {
			return nil, err
		}
	}

	return &Leaderboard{
		scoreboardFile:   scoreboardFile,
		existsStorageDir: true,
	}, nil
}

func (l *Leaderboard) update(player string, score int) ([]*Score, bool) {
	scores, ok := l.get()
	if !ok {
		return nil, false
	}

	scores = append(scores, &Score{player, score})
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})
	if len(scores) > MAX_FILE_LINES {
		scores = scores[:MAX_FILE_LINES]
	}

	rows := []string{}
	for _, s := range scores {
		rows = append(rows, s.player+"\t"+strconv.Itoa(s.score))
	}
	scoresData := strings.Join(rows, "\n")
	err := os.WriteFile(l.scoreboardFile, []byte(scoresData), 0644)
	if err != nil {
		return nil, false
	}

	return scores, true
}

func (l *Leaderboard) get() ([]*Score, bool) {
	if !l.existsStorageDir {
		return nil, false
	}

	data, err := os.ReadFile(l.scoreboardFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, true
		}
		return nil, false
	} else if len(data) == 0 {
		return []*Score{}, true
	}

	content := strings.Split(string(data), "\n")
	scores := []*Score{}
	for _, r := range content {
		row := strings.Split(r, "\t")
		player := row[0]
		if score, err := strconv.Atoi(row[1]); err == nil {
			scores = append(scores, &Score{player, score})
		}
	}

	return scores, true
}
