package gsnake

import (
	"sort"
	"sync"
	"time"
	"turutupa/gsnake/log"
)

const POINTS_PER_DEATH int = 20

var pos1 = [2]int{10, 10}
var pos2 = [2]int{10, COLS_MULTI - 10}
var pos3 = [2]int{ROWS_MULTI - 10, 10}
var pos4 = [2]int{ROWS_MULTI - 10, COLS_MULTI - 10}
var startingPositions = [][2]int{pos1, pos2, pos3, pos4}

type MultiGame struct {
	board   *Board
	players []*Player
	speed   int
}

func NewMultiGame(board *Board) *MultiGame {
	return &MultiGame{board, []*Player{}, ONLINE_SPEED}
}

func (g *MultiGame) AddPlayer(player *Player) {
	if player.snake == nil {
		sp := startingPositions[len(g.players)]
		player.WithSnake(sp[0], sp[1])
		log.Info("Created snake for player %s; id: %s", player.name, player.id)
	}
	player.snake.color = colors[len(g.players)+1]
	player.snake.arrowRenders = arrows[len(g.players)]
	g.players = append(g.players, player)
}

func (g *MultiGame) DefaultLayout(started bool) {
	for i, player := range g.players {
		player.screen.Clear()
		g.board.UpdateSnake(player.snake)
		if i == 1 || i == 3 {
			player.snake.Point(LEFT)
		}
	}

	for _, player := range g.players {
		player.screen.RenderBoardFrame(g.board)
		if len(g.players) < MAX_ROOM_SIZE && !started {
			player.screen.RenderWarning(g.board, "Waiting for players or game starts in 20s")
		}
		for _, p := range g.players {
			player.screen.RenderSnake(g.board, p.snake)
		}
	}
}

func (g *MultiGame) Countdown(label string) {
	for i := 5; i >= 0; i-- {
		var wg sync.WaitGroup
		for _, player := range g.players {
			wg.Add(1)
			go func(i int, player *Player) {
				defer wg.Done()
				player.screen.RenderCountdown(g.board, label, i)
				time.Sleep(time.Second)
			}(i, player)
		}
		wg.Wait()
	}
}

func (g *MultiGame) Run() {
	for {
		roundDeaths := 0
		for _, player := range g.players {
			if !player.isAlive {
				return
			}
			player.snake.Grow(1)
			player.snake.Move()
			if g.board.IntersectsMulti(player.snake) {
				player.isAlive = false
				roundDeaths++
			}
			g.board.UpdateSnake(player.snake)
		}

		for _, player := range g.players {
			for _, p := range g.players {
				if p.isAlive {
					player.screen.RenderSnake(g.board, p.snake)
				}
			}
			if player.isAlive {
				player.score = player.score + roundDeaths*POINTS_PER_DEATH
			}
		}

		if roundDeaths == len(g.players) {
			break
		}

		duration := time.Duration(g.speed) * time.Millisecond
		time.Sleep(duration)
	}
}

func (g *MultiGame) Stop() {}

func (g *MultiGame) RestartRound() {
	g.board.Restart()
	for i, p := range g.players {
		p.isAlive = true
		sp := startingPositions[i]
		p.snake.Restart(sp[0], sp[1])
		if i == 1 || i == 3 {
			p.snake.Point(LEFT)
		}
	}
}

type ByScore []*Player

func (p ByScore) Len() int           { return len(p) }
func (p ByScore) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByScore) Less(i, j int) bool { return p[i].score > p[j].score }

func (g *MultiGame) Leaderboard() {
	sort.Sort(ByScore(g.players))
	for _, p := range g.players {
		p.screen.Clear()
		p.screen.RenderMultiLeaderboard(g.board, g.players, g.HasWinner())
	}
}

func (g *MultiGame) HasWinner() bool {
	targetScore := 40
	// targetScore := len(g.players) * POINTS_PER_DEATH
	for _, p := range g.players {
		if p.score >= targetScore {
			return true
		}
	}
	return false
}
