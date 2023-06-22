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

func (l *Leaderboard) update(score int) ([]int, bool) {
	scores, ok := l.get()
	if !ok {
		return nil, false
	}

	scores = append(scores, score)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i] > scores[j]
	})
	if len(scores) > MAX_FILE_LINES {
		scores = scores[:MAX_FILE_LINES]
	}

	scoresData := strings.Join(intSliceToStringSlice(scores), "\n")
	err := os.WriteFile(l.scoreboardFile, []byte(scoresData), 0644)
	if err != nil {
		return nil, false
	}

	return scores, true
}

func (l *Leaderboard) get() ([]int, bool) {
	if !l.existsStorageDir {
		return nil, false
	}

	data, err := os.ReadFile(l.scoreboardFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, true
		}
		return nil, false
	}

	content := strings.Split(string(data), "\n")
	scores := []int{}
	for _, s := range content {
		if num, err := strconv.Atoi(s); err == nil {
			scores = append(scores, num)
		}
	}

	return scores, true
}

func intSliceToStringSlice(slice []int) []string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = strconv.Itoa(v)
	}
	return strSlice
}
