package gsnake

import (
	"math"
	"strconv"
	"time"

	"turutupa/gsnake/events"
)

type MenuOptionOffline int

// Main Menu options
const (
	EXIT        = "EXIT"
	LEADERBOARD = "LEADERBOARD"
	EASY        = "EASY"
	NORMAL      = "NORMAL"
	HARD        = "HARD"
	INSANITY    = "INSANITY"
)

var DIFFICULTIES = []string{EASY, NORMAL, HARD, INSANITY}
var MENU_OPTIONS = []string{EASY, NORMAL, HARD, INSANITY, EXIT}

// Snake available speeds
const (
	EASY_SPEED     = 50
	NORMAL_SPEED   = 40
	HARD_SPEED     = 30
	INSANITY_SPEED = 20
)

// Int value of arrow keys
const (
	ARROW_UP    int = 65
	ARROW_DOWN  int = 66
	ARROW_RIGHT int = 67
	ARROW_LEFT  int = 68
)

// Game State
type State int

const (
	MAIN_MENU              State = 1
	PLAYING                State = 2
	LEADERBOARD_MENU       State = 3
	LEADERBOARD_SUBMITTING State = 4
	FINISHED               State = 5
)

const MAX_PLAYER_LEN = 4

type Game struct {
	screen              *Screen
	eventPoller         events.EventPoller
	leaderboard         *Leaderboard
	fruit               *Fruit
	snake               *Snake
	player              string
	difficulty          string
	speed               int
	score               int
	state               State
	selectedMenuOption  int
	keypressCh          chan bool
	playerNameSubmitted chan bool
}

func NewGame(
	term events.EventPoller,
	screen *Screen,
	leaderboard *Leaderboard,
	fruit *Fruit,
	snake *Snake,
) *Game {
	game := &Game{
		screen:              screen,
		eventPoller:         term,
		leaderboard:         leaderboard,
		fruit:               fruit,
		snake:               snake,
		player:              "",
		speed:               0,
		score:               0,
		state:               MAIN_MENU,
		selectedMenuOption:  1, // defaults to NORMAL
		keypressCh:          make(chan bool),
		playerNameSubmitted: make(chan bool),
	}
	return game
}

func (g *Game) Run() {
	go g.executeUserInput()
	for {
		for g.state == MAIN_MENU {
			g.mainMenu()
			<-g.keypressCh // used for blocking
		}
		if g.state == FINISHED {
			g.Stop()
			return
		} else if g.state == LEADERBOARD_MENU {
			scores := g.leaderboard.get(g.difficulty)
			g.screen.renderScoreboard(g.difficulty, scores, nil)
			<-g.keypressCh
		} else if g.state == PLAYING {
			g.screen.clearTerminal()
			g.runGame()
		}
	}
}

func (g *Game) Stop() {
	select {
	case _, ok := <-g.keypressCh:
		if ok {
			close(g.keypressCh)
		}
	default:
	}
	g.eventPoller.Close()
	g.screen.clearTerminal()
	g.screen.showCursor()
}

func (g *Game) restart() {
	g.player = ""
	g.score = 0
	g.screen.restart()
	g.snake.restart(g.screen)
	g.fruit.new()
	g.state = MAIN_MENU
}

func (g *Game) mainMenu() {
	g.screen.clearTerminal()
	g.screen.renderMainMenu(g.selectedMenuOption)
}

func (g *Game) runGame() {
	var frameStart, frameTime uint32
	var frameEnd time.Time
	g.screen.init()
	for g.state == PLAYING {
		frameStart = uint32(time.Now().UnixNano() / int64(time.Millisecond)) // Current time in milliseconds
		g.screen.clear(g.fruit, g.snake.head, g.snake.tail, g.score)
		g.snake.move()
		if g.ateFruit() {
			g.score += 10
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
		g.screen.update(g.fruit, g.snake.head, g.score)
		g.screen.renderSnake(g.fruit, g.snake.head, g.snake.tail, g.score)
		if g.intersects() {
			time.Sleep(1 * time.Second)
			g.screen.clearTerminal()
			g.screen.printLogo()
			isHighScore := g.leaderboard.isHighScore(g.difficulty, g.score)
			scores := g.leaderboard.get(g.difficulty)
			if isHighScore {
				g.state = LEADERBOARD_SUBMITTING
				g.screen.renderScoreboard(g.difficulty, scores, &Score{g.player, g.score})
				isSubmitting := true
				for isSubmitting {
					select {
					case <-g.keypressCh:
						g.screen.renderScoreboard(g.difficulty, scores, &Score{g.player, g.score})
					case <-g.playerNameSubmitted:
						g.leaderboard.update(g.player, g.difficulty, g.score)
						isSubmitting = false
						break
					}
				}
			} else {
				g.state = LEADERBOARD_MENU
				g.screen.renderScoreboard(g.difficulty, scores, nil)
				<-g.keypressCh
			}
			g.restart()
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

func (g *Game) executeUserInput() {
	for {
		event, err := g.eventPoller.Poll()
		if err != nil {
			g.onExit()
			return
		}
		e := rune(event)
		if g.state == MAIN_MENU {
			g.userActionMainMenu(e)
		} else if g.state == PLAYING {
			g.userActionSnake(e)
		} else if g.state == LEADERBOARD_MENU {
			g.userActionLeaderboardMenu(e)
		} else if g.state == LEADERBOARD_SUBMITTING {
			g.userActionLeaderboardSubmitting(e)
		}
	}
}

func (g *Game) userActionMainMenu(event rune) {
	if g.isUp(event) {
		g.selectedMenuOption = int(math.Max(float64(0), float64(g.selectedMenuOption-1)))
	} else if g.isDown(event) {
		g.selectedMenuOption = int(math.Min(float64(len(MENU_OPTIONS)-1), float64(g.selectedMenuOption+1)))
	} else if g.isEnterKey(event) {
		selectedOpt := MENU_OPTIONS[g.selectedMenuOption]
		if selectedOpt == EXIT {
			g.onExit()
		} else if selectedOpt == LEADERBOARD {
			g.state = LEADERBOARD_MENU
		} else {
			if selectedOpt == EASY {
				g.speed = EASY_SPEED
				g.difficulty = EASY
			} else if selectedOpt == NORMAL {
				g.speed = NORMAL_SPEED
				g.difficulty = NORMAL
			} else if selectedOpt == HARD {
				g.speed = HARD_SPEED
				g.difficulty = HARD
			} else if selectedOpt == INSANITY {
				g.speed = INSANITY_SPEED
				g.difficulty = INSANITY
			}
			g.state = PLAYING
		}
	} else if event == 'q' {
		g.onExit()
	}
	g.keypressCh <- true
}

func (g *Game) userActionLeaderboardMenu(event rune) {
	if event == 'q' || g.isEnterKey(event) {
		g.onExit()
		g.keypressCh <- true
	}
}

func (g *Game) userActionLeaderboardSubmitting(event rune) {
	if isBackspaceOrDelete(event) {
		if len(g.player) > 0 {
			g.player = g.player[:len(g.player)-1]
		}
	}
	if isUserAcceptedChar(event) {
		if len(g.player) < MAX_PLAYER_LEN {
			g.player = g.player + string(event)
		}
	} else if g.isEnterKey(event) {
		g.playerNameSubmitted <- true
		return
	}
	g.keypressCh <- true
}

func (g *Game) userActionSnake(event rune) {
	pointing := g.snake.pointsTo()
	if g.isUp(event) {
		if pointing != DOWN {
			g.snake.point(UP)
		}
	} else if g.isDown(event) {
		if pointing != UP {
			g.snake.point(DOWN)
		}
	} else if g.isLeft(event) {
		if pointing != RIGHT {
			g.snake.point(LEFT)
		}
	} else if g.isRight(event) {
		if pointing != LEFT {
			g.snake.point(RIGHT)
		}
	} else if event == 'q' {
		g.onExit()
	}
}

func (g *Game) onExit() {
	if g.state == MAIN_MENU {
		g.state = FINISHED
		return
	}
	g.restart()
}

func isBackspaceOrDelete(r rune) bool {
	return r == '\b' || r == '\u007F'
}

func isUserAcceptedChar(r rune) bool {
	return byte(r) >= 33 && byte(r) <= 126
}

// Accepted keys for up/down/left and right are
// - wasd
// - hjkl
// - arrow keys
func (g *Game) isUp(event rune) bool {
	return event == 'w' || int(event) == ARROW_UP || event == 'k'
}

func (g *Game) isDown(event rune) bool {
	return event == 's' || int(event) == ARROW_DOWN || event == 'j'
}

func (g *Game) isLeft(event rune) bool {
	return event == 'a' || int(event) == ARROW_LEFT || event == 'h'
}

func (g *Game) isRight(event rune) bool {
	return event == 'd' || int(event) == ARROW_RIGHT || event == 'l'
}

// enter keys are
// - enter
// - spacebar
// - \r which I'm not sure which key is that tbh
func (g *Game) isEnterKey(input rune) bool {
	in := byte(input)
	enterKeys := [2]byte{'\n', '\r'} // Byte representations of "enter" keys
	for _, key := range enterKeys {
		if in == key {
			return true
		}
	}
	if int(in) == 32 { // space
		return true
	}
	return false
}

func toInt(str string) int {
	s, _ := strconv.Atoi(str)
	return s
}
