package gsnake

type Pointing rune

const (
	UP    Pointing = 'u'
	DOWN  Pointing = 'd'
	RIGHT Pointing = 'r'
	LEFT  Pointing = 'l'
)

const (
	arrowUp0    Pointing = '▲'
	arrowDown0  Pointing = '▼'
	arrowRight0 Pointing = '►'
	arrowLeft0  Pointing = '◄'
)

const (
	arrowUp1    Pointing = '△'
	arrowDown1  Pointing = '▽'
	arrowRight1 Pointing = '▷'
	arrowLeft1  Pointing = '◁'
)

const (
	arrowUp2    Pointing = '⇡'
	arrowDown2  Pointing = '⇣'
	arrowRight2 Pointing = '⇢'
	arrowLeft2  Pointing = '⇠'
)

// Set 5
const (
	arrowUp3    Pointing = '⮝'
	arrowDown3  Pointing = '⮟'
	arrowRight3 Pointing = '⮞'
	arrowLeft3  Pointing = '⮜'
)

var arrows = []map[rune]Pointing{
	{'u': arrowUp0, 'd': arrowDown0, 'l': arrowLeft0, 'r': arrowRight0},
	{'u': arrowUp1, 'd': arrowDown1, 'l': arrowLeft1, 'r': arrowRight1},
	{'u': arrowUp2, 'd': arrowDown2, 'l': arrowLeft2, 'r': arrowRight2},
	{'u': arrowUp3, 'd': arrowDown3, 'l': arrowLeft3, 'r': arrowRight3},
}

type Color = string

const (
	white   = "\033[37m"
	red     = "\033[31m"
	green   = "\033[32m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	reset   = "\033[0m"
)

var colors = []Color{white, red, green, blue, magenta}

type Snake struct {
	head         *Node
	tail         *Node
	color        Color
	arrowRenders map[rune]Pointing
}

type Node struct {
	row         int
	col         int
	pointing    Pointing
	tmpPointing Pointing
	prev        *Node
	next        *Node
	render      rune
	validated   bool
}

func NewSnake(row int, col int) *Snake {
	snake := &Snake{}
	snake.arrowRenders = arrows[0] // default
	snake.color = colors[0]        // default
	snake.head = &Node{
		row:         row,
		col:         col,
		pointing:    RIGHT,
		tmpPointing: RIGHT,
		prev:        nil,
		next:        nil,
		render:      HORIZONTAL,
		validated:   true,
	}
	snake.tail = snake.head
	return snake
}

func (s *Snake) Restart(row int, col int) {
	newSnake := NewSnake(row, col)
	s.head = newSnake.head
	s.tail = newSnake.tail
}

func (s *Snake) PointsTo() Pointing {
	return s.head.pointing
}

func (s *Snake) Point(point Pointing) {
	s.head.tmpPointing = point
	switch point {
	case UP:
		s.head.render = rune(s.arrowRenders['u'])
	case DOWN:
		s.head.render = rune(s.arrowRenders['d'])
	case LEFT:
		s.head.render = rune(s.arrowRenders['l'])
	case RIGHT:
		s.head.render = rune(s.arrowRenders['r'])
	}
}

func (s *Snake) Move() {
	node := s.head
	x_prev := node.row
	y_prev := node.col
	pointing_prev := node.tmpPointing

	switch node.tmpPointing {
	case UP:
		node.row = node.row - 1
	case RIGHT:
		node.col = node.col + 1
	case LEFT:
		node.col = node.col - 1
	case DOWN:
		node.row = node.row + 1
	}

	node.pointing = node.tmpPointing
	node = node.next
	for node != nil {
		pointing_tmp := node.tmpPointing
		x_tmp := node.row
		y_tmp := node.col
		node.row = x_prev
		node.col = y_prev
		node.pointing = pointing_prev
		node.tmpPointing = pointing_prev
		x_prev = x_tmp
		y_prev = y_tmp
		pointing_prev = pointing_tmp
		if !node.validated {
			node.validated = true
			break
		}
		node = node.next
	}
}

func (s *Snake) Grow(size int) {
	for i := 0; i < size; i++ {
		s.append()
	}
}

func (s *Snake) append() {
	var x int
	var y int
	if s.tail.pointing == UP {
		x = s.tail.row + 1
		y = s.tail.col
	} else if s.tail.pointing == RIGHT {
		x = s.tail.row
		y = s.tail.col - 1
	} else if s.tail.pointing == DOWN {
		x = s.tail.row - 1
		y = s.tail.col
	} else {
		x = s.tail.row
		y = s.tail.col + 1
	}
	s.tail.next = &Node{
		row:       x,
		col:       y,
		pointing:  s.tail.pointing,
		prev:      s.tail,
		next:      nil,
		validated: false,
	}
	s.tail = s.tail.next
}

func snakeStrategy(snake *Snake, event rune) {
	pointing := snake.PointsTo()
	if isUp(event) {
		if pointing != DOWN {
			snake.Point(UP)
		}
	} else if isDown(event) {
		if pointing != UP {
			snake.Point(DOWN)
		}
	} else if isLeft(event) {
		if pointing != RIGHT {
			snake.Point(LEFT)
		}
	} else if isRight(event) {
		if pointing != LEFT {
			snake.Point(RIGHT)
		}
	}
}
