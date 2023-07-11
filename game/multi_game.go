package gsnake

import "time"

type MultiGame struct {
	board   *Board
	screen  *Screen
	players []*Player
	snakes  []*Snake
}

func (g *MultiGame) Setup() {
	g.screen.RenderBoard(g.board)
	for {
		for _, snake := range g.snakes {
			g.screen.Remove(g.board, snake.head)
			snake.Move()
			g.board.UpdateSnake(snake.head)
		}
	}
}

func (g *MultiGame) Countdown() {
	for i := 1; i <= 5; i++ {
		for _, snake := range g.snakes {
			g.screen.RenderSnake(g.board, snake.head)
		}
		g.screen.RenderCountdown(i)
		time.Sleep(1 * time.Second)
	}
}

func (g *MultiGame) Run() {}

func (g *MultiGame) Stop() {}

func (g *MultiGame) Restart() {}

func (g *MultiGame) InterRound() {}

func (g *MultiGame) Leaderboard() {}

func (g *MultiGame) HasWinner() bool { return false }

func (g *MultiGame) Strategy(event rune) {}
