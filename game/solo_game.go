package gsnake

import (
	"time"
)

// Snake available speeds
const (
	EASY_SPEED     = 50
	NORMAL_SPEED   = 40
	HARD_SPEED     = 30
	INSANITY_SPEED = 20
)

const MAX_PLAYER_NAME_LEN = 4

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
	g.screen.RenderBoard(g.board)
	g.screen.RenderScore(g.board, 0)
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
		g.board.UpdateSnake(g.snake.head)

		// render board
		g.screen.RenderScore(g.board, g.player.score)
		g.screen.RenderFruit(g.board, g.fruit)
		g.screen.RenderSnake(g.board, g.snake.head)
		if g.board.Intersects(g.snake) {
			time.Sleep(1 * time.Second)
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
		g.screen.RenderLeaderboard(g.board, g.difficulty, scores, g.player)
		isSubmitting := true
		for isSubmitting {
			select {
			case <-g.keypressCh:
				g.screen.RenderLeaderboard(g.board, g.difficulty, scores, g.player)
			case <-g.playerNameSubmitted:
				g.leaderboard.Update(g.difficulty, g.player)
				isSubmitting = false
				break
			}
		}
	} else {
		g.state = LEADERBOARD_MENU
		g.screen.RenderLeaderboard(g.board, g.difficulty, scores, nil)
		<-g.keypressCh
	}
}

func (g *SoloGame) Stop() {
	g.state = FINISHED
}

func (g *SoloGame) Restart() {
	g.player = NewPlayer("")
	g.screen.Restart()
	g.snake.Restart(g.board)
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
	x := g.snake.head.x
	y := g.snake.head.y
	return x == g.fruit.x && y == g.fruit.y
}

func (g *SoloGame) Strategy(event rune) {
	switch g.state {
	case PLAYING:
		g.userActionSnake(event)
	case LEADERBOARD_MENU:
		g.userActionLeaderboardMenu(event)
	case LEADERBOARD_SUBMITTING:
		g.userActionLeaderboardSubmitting(event)
	}
}

func (g *SoloGame) userActionLeaderboardMenu(event rune) {
	if event == 'q' || isEnterKey(event) {
		g.state = FINISHED
		g.keypressCh <- true
	}
}

func (g *SoloGame) userActionLeaderboardSubmitting(event rune) {
	name, done := HandleUserInputForm(g.player.name, event)
	if done {
		g.playerNameSubmitted <- true
		return
	}
	g.player.name = name
	g.keypressCh <- true
}

func (g *SoloGame) userActionSnake(event rune) {
	pointing := g.snake.PointsTo()
	if isUp(event) {
		if pointing != DOWN {
			g.snake.Point(UP)
		}
	} else if isDown(event) {
		if pointing != UP {
			g.snake.Point(DOWN)
		}
	} else if isLeft(event) {
		if pointing != RIGHT {
			g.snake.Point(LEFT)
		}
	} else if isRight(event) {
		if pointing != LEFT {
			g.snake.Point(RIGHT)
		}
	} else if event == 'q' {
		g.player.score = 0 // he quit so player shouldn't set new high score
		g.state = FINISHED
	}
}
