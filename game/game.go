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

const MAX_PLAYER_LEN = 4

type GameState int

const (
	PLAYING                GameState = 1
	LEADERBOARD_MENU       GameState = 2
	LEADERBOARD_SUBMITTING GameState = 3
	FINISHED               GameState = 4
)

type Game struct {
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

func NewGame(
	screen *Screen,
	leaderboard *Leaderboard,
	fruit *Fruit,
	snake *Snake,
) *Game {
	game := &Game{
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

func (g *Game) Run() {
	var frameStart, frameTime uint32
	var frameEnd time.Time
	g.screen.renderBoard()
	for g.state == PLAYING {
		frameStart = uint32(time.Now().UnixNano() / int64(time.Millisecond)) // Current time in milliseconds
		g.screen.remove(g.snake.head, g.snake.tail)
		g.snake.move()
		if g.ateFruit() {
			g.player.score += 10
			if g.speed == EASY_SPEED || g.speed == NORMAL_SPEED {
				g.snake.append()
				g.snake.append()
			} else {
				g.snake.append()
				g.snake.append()
				g.snake.append()
			}
			g.fruit.new()
		}
		g.screen.update(g.fruit, g.snake.head, g.player.score)
		g.screen.renderSnake(g.fruit, g.snake.head, g.snake.tail, g.player.score)
		if g.intersects() {
			time.Sleep(1 * time.Second)
			g.screen.clear()
			g.screen.printLogo()
			isHighScore := g.leaderboard.isHighScore(g.difficulty, g.player.score)
			scores := g.leaderboard.get(g.difficulty)
			if isHighScore {
				g.state = LEADERBOARD_SUBMITTING
				g.screen.renderScoreboard(g.difficulty, scores, g.player)
				isSubmitting := true
				for isSubmitting {
					select {
					case <-g.keypressCh:
						g.screen.renderScoreboard(g.difficulty, scores, g.player)
					case <-g.playerNameSubmitted:
						g.leaderboard.update(g.difficulty, g.player)
						isSubmitting = false
						break
					}
				}
			} else {
				g.state = LEADERBOARD_MENU
				g.screen.renderScoreboard(g.difficulty, scores, nil)
				<-g.keypressCh
			}
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

func (g *Game) Stop() {
	g.state = FINISHED
}

func (g *Game) Restart() {
	g.player = NewPlayer("")
	g.screen.restart()
	g.snake.restart(g.screen)
	g.fruit.new()
	g.state = PLAYING
}

func (g *Game) setDifficulty(diff string) {
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

func (g *Game) ateFruit() bool {
	x := g.snake.head.x
	y := g.snake.head.y
	return x == g.fruit.x && y == g.fruit.y
}

func (g *Game) intersects() bool {
	head := g.snake.head
	x := head.x
	y := head.y
	if x == 0 || x == g.screen.rows-1 || y == 0 || y == g.screen.cols-1 {
		return true
	}

	node := head.next
	for node != nil && node.validated {
		if x == node.x && y == node.y {
			return true
		}
		node = node.next
	}
	return false
}

func (g *Game) strategy(event rune) {
	switch g.state {
	case PLAYING:
		g.userActionSnake(event)
	case LEADERBOARD_MENU:
		g.userActionLeaderboardMenu(event)
	case LEADERBOARD_SUBMITTING:
		g.userActionLeaderboardSubmitting(event)
	}
}

func (g *Game) userActionLeaderboardMenu(event rune) {
	if event == 'q' || isEnterKey(event) {
		g.state = FINISHED
		g.keypressCh <- true
	}
}

func (g *Game) userActionLeaderboardSubmitting(event rune) {
	if isBackspaceOrDelete(event) {
		if len(g.player.name) > 0 {
			g.player.name = g.player.name[:len(g.player.name)-1]
		}
	}
	if isUserAcceptedChar(event) {
		if len(g.player.name) < MAX_PLAYER_LEN {
			g.player.name = g.player.name + string(event)
		}
	} else if isEnterKey(event) {
		g.playerNameSubmitted <- true
		return
	}
	g.keypressCh <- true
}

func (g *Game) userActionSnake(event rune) {
	pointing := g.snake.pointsTo()
	if isUp(event) {
		if pointing != DOWN {
			g.snake.point(UP)
		}
	} else if isDown(event) {
		if pointing != UP {
			g.snake.point(DOWN)
		}
	} else if isLeft(event) {
		if pointing != RIGHT {
			g.snake.point(LEFT)
		}
	} else if isRight(event) {
		if pointing != LEFT {
			g.snake.point(RIGHT)
		}
	} else if event == 'q' {
		g.state = FINISHED
	}
}
