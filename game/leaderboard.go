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

type Leaderboard struct {
	enabled     bool
	leaderboard map[string][]*Player
}

func NewLeaderboard() *Leaderboard {
	l := &Leaderboard{true, make(map[string][]*Player)}
	for _, diff := range DIFFICULTIES {
		filename := l.getFilename(diff)
		_, ok := fsutil.NewCfgFile(filename)
		if !ok {
			l.enabled = false
		}
		if scores, ok := l.getPersistedLeaderboard(filename); ok {
			l.leaderboard[diff] = scores
		} else {
			log.Error("Something went wrong retrieving local leaderboard")
		}
	}
	return l
}

func (l *Leaderboard) IsHighScore(difficulty string, score int) bool {
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

func (l *Leaderboard) Get(difficulty string) []*Player {
	return l.leaderboard[difficulty]
}

func (l *Leaderboard) Update(difficulty string, player *Player) ([]*Player, bool) {
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
	players, exists := l.leaderboard[difficulty]
	if !exists {
		players = []*Player{}
	}
	players = append(players, player)
	sort.Slice(players, func(i, j int) bool {
		return players[i].score > players[j].score
	})
	if len(players) > MAX_SCORES_STORED {
		players = players[:MAX_SCORES_STORED]
	}
	l.leaderboard[difficulty] = players

	// persist
	rows := []string{}
	for _, s := range players {
		rows = append(rows, s.name+"\t"+strconv.Itoa(s.score))
	}
	playersData := strings.Join(rows, "\n")
	err := fsutil.WriteFile(l.getFilename(difficulty), []byte(playersData))
	if err != nil {
		return nil, false
	}

	return players, true
}

func (l *Leaderboard) getPersistedLeaderboard(filename string) ([]*Player, bool) {
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
		return []*Player{}, true
	}

	content := strings.Split(string(data), "\n")
	scores := []*Player{}
	for _, r := range content {
		row := strings.Split(r, "\t")
		if len(row) != 2 {
			return nil, false
		}
		player := row[0]
		if score, err := strconv.Atoi(row[1]); err == nil {
			scores = append(scores, NewPlayer(player).WithScore(score))
		} else {
			return nil, false
		}
	}

	return scores, true
}

func (l *Leaderboard) getFilename(difficulty string) string {
	return FILENAME_PREFIX + strings.ToLower(difficulty)
}
