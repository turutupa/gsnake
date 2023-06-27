package gsnake

import (
	"os"
	"sort"
	"strconv"
	"strings"
	"turutupa/gsnake/fsutil"
	"turutupa/gsnake/log"
)

const MAX_SCORES_STORED = 50
const FILENAME_PREFIX = "leaderboard."

type Score struct {
	player string
	score  int
}

type Leaderboard struct {
	enabled     bool
	leaderboard map[string][]*Score
}

func NewLeaderboard() *Leaderboard {
	l := &Leaderboard{true, make(map[string][]*Score)}
	for _, diff := range DIFFICULTIES {
		filename := l.getFilename(diff)
		_, ok := fsutil.NewCfgFile(filename)
		if !ok {
			l.enabled = false
		}
		if scores, ok := l.getPersistedScores(filename); ok {
			l.leaderboard[diff] = scores
		} else {
			log.Error("Something went wrong retrieving local scores", nil)
		}
	}
	return l
}

func (l *Leaderboard) isHighScore(difficulty string, score int) bool {
	if score == 0 {
		return false
	}
	if len(l.leaderboard[difficulty]) < 5 {
		return true
	}
	for i, s := range l.leaderboard[difficulty] {
		if i >= 5 {
			return false
		}
		if score > s.score {
			return true
		}
	}
	return false
}

func (l *Leaderboard) get(difficulty string) []*Score {
	return l.leaderboard[difficulty]
}

func (l *Leaderboard) update(player string, difficulty string, score int) ([]*Score, bool) {
	if !l.enabled {
		return nil, false
	}

	existsDiff := false
	for _, diff := range DIFFICULTIES {
		if diff == difficulty {
			existsDiff = true
		}
	}

	if !existsDiff {
		log.Error("Leaderboard trying to update a difficulty that doesn't exist.", nil)
		return nil, false
	}

	scores, exists := l.leaderboard[difficulty]
	if !exists {
		scores = []*Score{}
	}

	scores = append(scores, &Score{player, score})
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	if len(scores) > MAX_SCORES_STORED {
		scores = scores[:MAX_SCORES_STORED]
	}

	l.leaderboard[difficulty] = scores

	rows := []string{}
	for _, s := range scores {
		rows = append(rows, s.player+"\t"+strconv.Itoa(s.score))
	}
	scoresData := strings.Join(rows, "\n")
	err := fsutil.WriteFile(l.getFilename(difficulty), []byte(scoresData))
	if err != nil {
		return nil, false
	}

	return l.leaderboard[difficulty], true
}

func (l *Leaderboard) getPersistedScores(filename string) ([]*Score, bool) {
	if !l.enabled {
		return nil, false
	}

	data, err := fsutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, true
		}
		return nil, false
	} else if len(data) == 0 {
		// file is empty
		return []*Score{}, true
	}

	content := strings.Split(string(data), "\n")
	scores := []*Score{}
	for _, r := range content {
		row := strings.Split(r, "\t")
		if len(row) != 2 {
			return nil, false
		}
		player := row[0]
		if score, err := strconv.Atoi(row[1]); err == nil {
			scores = append(scores, &Score{player, score})
		} else {
			return nil, false
		}
	}

	return scores, true
}

func (l *Leaderboard) getFilename(difficulty string) string {
	return FILENAME_PREFIX + strings.ToLower(difficulty)
}
