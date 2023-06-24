package gsnake

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"turutupa/gsnake/fsutil"
)

const FILENAME = "scoreboard"
const MAX_FILE_LINES = 100

type Score struct {
	player string
	score  int
}

type Leaderboard struct {
	enabled bool
}

func NewLeaderboard() *Leaderboard {
	_, ok := fsutil.NewCfgFile(FILENAME)
	return &Leaderboard{ok}
}

func (l *Leaderboard) isHighScore(score int) (bool, bool) {
	scores, ok := l.get()
	if !ok {
		return false, false
	}
	if len(scores) < 5 {
		return true, true
	}
	for i, s := range scores {
		if i >= 5 {
			return false, true
		}
		if score > s.score {
			return true, true
		}
	}
	return false, true
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
	err := fsutil.WriteFile(FILENAME, []byte(scoresData))
	if err != nil {
		return nil, false
	}

	return scores, true
}

func (l *Leaderboard) get() ([]*Score, bool) {
	if !l.enabled {
		return nil, false
	}

	data, err := fsutil.ReadFile(FILENAME)
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
