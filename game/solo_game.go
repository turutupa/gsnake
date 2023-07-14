package gsnake

import (
	"time"
)

type GameState int

const (
	PLAYING                GameState = 1
	LEADERBOARD_MENU       GameState = 2
	LEADERBOARD_SUBMITTING GameState = 3
	FINISHED               GameState = 4
)

type SoloGame struct {
	board               *Board
	screen              *Screen
	leaderboard         *Leaderboard
	fruit               *Fruit
	snake               *Snake
	player              *Player
	difficulty          string
	speed               int
	state               GameState
	keypressCh          chan bool
	playerNameSubmitted chan bool
}

func NewSoloGame(
	board *Board,
	screen *Screen,
	leaderboard *Leaderboard,
	fruit *Fruit,
	snake *Snake,
) *SoloGame {
	game := &SoloGame{
		board:               board,
		screen:              screen,
		leaderboard:         leaderboard,
		fruit:               fruit,
		snake:               snake,
		player:              NewPlayer(""),
		speed:               0,
		difficulty:          NORMAL,
		state:               PLAYING,
		keypressCh:          make(chan bool),
		playerNameSubmitted: make(chan bool),
	}
	return game
}

func (g *SoloGame) Run() {
	var frameStart, frameTime uint32
	var frameEnd time.Time
	g.screen.RenderBoardFrame(g.board)
	g.screen.RenderScore(g.board, 0)
	g.snake.Grow(6) // snake by default has length 1 so we increase to 7
	for g.state == PLAYING {
		frameStart = uint32(time.Now().UnixNano() / int64(time.Millisecond)) // Current time in milliseconds
		g.screen.Remove(g.board, g.snake.head)
		g.snake.Move()
		if g.ateFruit() {
			switch g.speed {
			case EASY_SPEED:
				g.snake.Grow(2)
				g.player.score += 10
			case NORMAL_SPEED:
				g.snake.Grow(2)
				g.player.score += 10
			case HARD_SPEED:
				g.snake.Grow(3)
				g.player.score += 10
			case INSANITY_SPEED:
				g.snake.Grow(3)
				g.player.score += 15
			}
			g.fruit.New()
		}
		// update internal board
		g.board.UpdateFruit(g.fruit)
		g.board.UpdateLeaderboard(g.player.score)
		g.board.UpdateSnake(g.snake)

		// render board
		g.screen.RenderScore(g.board, g.player.score)
		g.screen.RenderFruit(g.board, g.fruit)
		g.screen.RenderSnake(g.board, g.snake)

		// game over!
		if g.board.Intersects(g.snake) {
			time.Sleep(time.Second)
			return
		}

		//  adding some extra time when going vertical because it feels faster
		//  than when going horizontally due to font width/height.
		//  Will solve that when we make it double columned
		frameEnd = time.Now()
		frameTime = uint32(frameEnd.UnixNano()/int64(time.Millisecond)) - frameStart
		frameDelay := uint32(g.speed)
		sleepDuration := time.Duration(frameDelay-frameTime) * time.Millisecond
		if g.snake.head.pointing == UP || g.snake.head.pointing == DOWN {
			if frameDelay > frameTime {
				time.Sleep(sleepDuration * 2)
			}
		} else {
			if frameDelay > frameTime {
				time.Sleep(sleepDuration)
			}
		}
	}
}

func (g *SoloGame) Leaderboard() {
	g.screen.Clear()
	g.screen.PrintLogo()
	isHighScore := g.leaderboard.IsHighScore(g.difficulty, g.player.score)
	scores := g.leaderboard.Get(g.difficulty)
	if isHighScore {
		g.state = LEADERBOARD_SUBMITTING
		g.screen.RenderSingleLeaderboard(g.board, g.difficulty, scores, g.player)
		isSubmitting := true
		for isSubmitting {
			select {
			case <-g.keypressCh:
				g.screen.RenderSingleLeaderboard(g.board, g.difficulty, scores, g.player)
			case <-g.playerNameSubmitted:
				g.leaderboard.Update(g.difficulty, g.player)
				isSubmitting = false
				break
			}
		}
	} else {
		g.state = LEADERBOARD_MENU
		g.screen.RenderSingleLeaderboard(g.board, g.difficulty, scores, nil)
		<-g.keypressCh
	}
}

func (g *SoloGame) Stop() {
	g.state = FINISHED
}

func (g *SoloGame) Restart() {
	g.player = NewPlayer("")
	g.screen.Restart()
	g.snake.Restart(g.board.rows/2, g.board.cols/5)
	g.fruit.New()
	g.state = PLAYING
}

func (g *SoloGame) SetDifficulty(diff string) {
	switch diff {
	case EASY:
		g.speed = EASY_SPEED
		g.difficulty = EASY
	case NORMAL:
		g.speed = NORMAL_SPEED
		g.difficulty = NORMAL
	case HARD:
		g.speed = HARD_SPEED
		g.difficulty = HARD
	case INSANITY:
		g.speed = INSANITY_SPEED
		g.difficulty = INSANITY
	}
}

func (g *SoloGame) ateFruit() bool {
	x := g.snake.head.row
	y := g.snake.head.col
	return x == g.fruit.x && y == g.fruit.y
}

func (g *SoloGame) Strategy(event rune) {
	switch g.state {
	case PLAYING:
		if event == 'q' {
			g.player.score = 0 // he quit so player shouldn't set new high score
			g.state = FINISHED
			return
		}
		snakeStrategy(g.snake, event)
	case LEADERBOARD_MENU:
		g.singleLeaderboardStrategy(event)
	case LEADERBOARD_SUBMITTING:
		g.singleLeaderboardSubmitStrategy(event)
	}
}

func (g *SoloGame) singleLeaderboardStrategy(event rune) {
	if event == 'q' || isEnterKey(event) {
		g.state = FINISHED
		g.keypressCh <- true
	}
}

func (g *SoloGame) singleLeaderboardSubmitStrategy(event rune) {
	name, done := HandleUserInputForm(g.player.name, event)
	if done {
		g.playerNameSubmitted <- true
		return
	}
	g.player.name = name
	g.keypressCh <- true
}
