package gsnake

import (
	"math"
	"time"

	"turutupa/gsnake/events"
)

type Speed int

const (
	EXIT     Speed = 0
	EASY     Speed = 50
	NORMAL   Speed = 40
	HARD     Speed = 30
	INSANITY Speed = 20
)

const ARROW_UP int = 65
const ARROW_DOWN int = 66
const ARROW_RIGHT int = 67
const ARROW_LEFT int = 68

var MENU_OPTIONS = []Speed{EASY, NORMAL, HARD, INSANITY, EXIT}

type State int

const (
	MAIN_MENU  State = 1
	PLAYING    State = 2
	SCOREBOARD State = 3
)

type Game struct {
	screen             *Screen
	eventPoller        events.EventPoller
	scoreboard         *Scoreboard
	fruit              *Fruit
	snake              *Snake
	speed              int
	running            bool
	score              int
	state              State
	selectedMenuOption int
	selectChan         chan bool
}

func NewGame(
	term events.EventPoller,
	screen *Screen,
	scoreboard *Scoreboard,
	fruit *Fruit,
	snake *Snake,
) *Game {
	game := &Game{
		screen:             screen,
		eventPoller:        term,
		scoreboard:         scoreboard,
		fruit:              fruit,
		snake:              snake,
		speed:              0,
		running:            true,
		score:              0,
		state:              MAIN_MENU,
		selectedMenuOption: 2,
		selectChan:         make(chan bool),
	}
	return game
}

func (g *Game) Run() {
	go g.executeUserInput()
	for g.running {
		for g.state == MAIN_MENU {
			g.mainMenu()
			<-g.selectChan // used for blocking
			if !g.running {
				break
			}
		}
		g.runGame()
	}

	g.screen.clearTerminal()
}

func (g *Game) Quit() {
	g.running = false
	g.eventPoller.Close()
	g.selectChan <- true
}

func (g *Game) restart() {
	g.score = 0
	g.screen.restart()
	g.snake.restart(g.screen)
	g.state = MAIN_MENU
}

func (g *Game) mainMenu() {
	g.screen.clearTerminal()
	g.screen.renderMainMenu(g.selectedMenuOption)
}

func (g *Game) runGame() {
	g.screen.clearTerminal()
	g.screen.init()
	for g.state == PLAYING {
		g.screen.clear(g.fruit, g.snake.head, g.snake.tail, g.score)
		g.snake.move()
		if g.ateFruit() {
			g.score += 10
			if g.speed == int(EASY) || g.speed == int(NORMAL) {
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
			scores, ok := g.scoreboard.update(g.score)
			time.Sleep(1 * time.Second)
			g.screen.clearTerminal()
			g.screen.GameOver()
			if ok {
				g.screen.renderScoreboard(scores)
			}
			g.state = SCOREBOARD
			<-g.selectChan
			g.restart()
			return
		}
		// adding some extra time when going vertical because it feels faster
		// than when going horizontally due to font width/height
		duration := time.Duration(g.speed) * time.Millisecond
		if g.snake.head.pointing == UP || g.snake.head.pointing == DOWN {
			time.Sleep(duration + (duration * 3 / 4))
		} else {
			time.Sleep(duration)
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
	for g.running {
		event := rune(g.eventPoller.Poll())
		if g.state == MAIN_MENU {
			g.userActionMainMenu(event)
		} else if g.state == PLAYING {
			g.userActionSnake(event)
		} else if g.state == SCOREBOARD {
			g.userActionScoreboard(event)
		}
	}
}

func (g *Game) userActionMainMenu(event rune) {
	if event == 'w' || event == 'k' || int(event) == ARROW_UP {
		g.selectedMenuOption = int(math.Max(float64(0), float64(g.selectedMenuOption-1)))
	} else if event == 's' || event == 'j' || int(event) == ARROW_DOWN {
		g.selectedMenuOption = int(math.Min(float64(len(MENU_OPTIONS)-1), float64(g.selectedMenuOption+1)))
	} else if g.isEnterKey(event) {
		if g.selectedMenuOption == len(MENU_OPTIONS)-1 {
			g.onExit()
			g.selectChan <- true
			return
		}
		g.speed = int(MENU_OPTIONS[g.selectedMenuOption])
		g.state = PLAYING
	} else if event == 'q' {
		g.onExit()
	}
	g.selectChan <- true
}

func (g *Game) userActionSnake(event rune) {
	pointing := g.snake.pointsTo()
	if event == 'w' || int(event) == ARROW_UP {
		if pointing != DOWN {
			g.snake.point(UP)
		}
	} else if event == 's' || int(event) == ARROW_DOWN {
		if pointing != UP {
			g.snake.point(DOWN)
		}
	} else if event == 'a' || int(event) == ARROW_LEFT {
		if pointing != RIGHT {
			g.snake.point(LEFT)
		}
	} else if event == 'd' || int(event) == ARROW_RIGHT {
		if pointing != LEFT {
			g.snake.point(RIGHT)
		}
	} else if event == 'q' {
		g.onExit()
	}
}

func (g *Game) userActionScoreboard(event rune) {
	if event == '\n' || event == 'q' || g.isEnterKey(event) {
		g.selectChan <- true
	}
}

func (g *Game) isEnterKey(input rune) bool {
	in := byte(input)
	enterKeys := [2]byte{'\n', '\r'} // Byte representations of "enter" keys
	for _, key := range enterKeys {
		if in == key {
			return true
		}
	}
	return false
}

func (g *Game) onExit() {
	if g.state == MAIN_MENU {
		g.screen.clearTerminal()
		g.Quit()
		return
	}
	g.restart()
}
