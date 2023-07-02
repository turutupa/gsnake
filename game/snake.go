package gsnake

type Pointing rune

// const (
// 	UP    Pointing = '▲'
// 	DOWN  Pointing = '▼'
// 	RIGHT Pointing = '►'
// 	LEFT  Pointing = '◄'
// )

// Set 1
const (
	UP    Pointing = '△'
	DOWN  Pointing = '▽'
	RIGHT Pointing = '▷'
	LEFT  Pointing = '◁'
)

// Set 2
const (
	UP2    Pointing = '⇧'
	DOWN2  Pointing = '⇩'
	RIGHT2 Pointing = '⇨'
	LEFT2  Pointing = '⇦'
)

// Set 3
const (
	UP3    Pointing = '⬆'
	DOWN3  Pointing = '⬇'
	RIGHT3 Pointing = '➡'
	LEFT3  Pointing = '⬅'
)

// Set 4
const (
	UP4    Pointing = '⇡'
	DOWN4  Pointing = '⇣'
	RIGHT4 Pointing = '⇢'
	LEFT4  Pointing = '⇠'
)

// Set 5
const (
	UP5    Pointing = '⮝'
	DOWN5  Pointing = '⮟'
	RIGHT5 Pointing = '⮞'
	LEFT5  Pointing = '⮜'
)

type Snake struct {
	head *Node
	tail *Node
}

type Node struct {
	x           int
	y           int
	pointing    Pointing
	tmpPointing Pointing
	prev        *Node
	next        *Node
	render      rune
	validated   bool
}

func NewSnake(board *Board) *Snake {
	snake := &Snake{}
	snake.head = &Node{
		x:           board.rows / 2,
		y:           board.cols / 5,
		pointing:    RIGHT,
		tmpPointing: RIGHT,
		prev:        nil,
		next:        nil,
		render:      HORIZONTAL,
		validated:   true,
	}
	node := snake.head
	for i := 0; i < 6; i++ { // initial length of 7
		node.next = &Node{
			x:           node.x,
			y:           node.y - 1,
			pointing:    RIGHT,
			tmpPointing: RIGHT,
			prev:        node,
			next:        nil,
			render:      HORIZONTAL,
			validated:   true,
		}
		node = node.next
	}
	snake.tail = node
	return snake
}

func (s *Snake) Restart(board *Board) {
	newSnake := NewSnake(board)
	s.head = newSnake.head
	s.tail = newSnake.tail
}

func (s *Snake) PointsTo() Pointing {
	return s.head.pointing
}

func (s *Snake) Point(point Pointing) {
	s.head.tmpPointing = point
}

func (s *Snake) Move() {
	node := s.head
	x_prev := node.x
	y_prev := node.y
	pointing_prev := node.tmpPointing
	if node.tmpPointing == UP {
		node.x = node.x - 1
	} else if node.tmpPointing == RIGHT {
		node.y = node.y + 1
	} else if node.tmpPointing == LEFT {
		node.y = node.y - 1
	} else {
		node.x = node.x + 1
	}
	node.pointing = node.tmpPointing
	node = node.next
	for node != nil {
		pointing_tmp := node.tmpPointing
		x_tmp := node.x
		y_tmp := node.y
		node.x = x_prev
		node.y = y_prev
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
		x = s.tail.x + 1
		y = s.tail.y
	} else if s.tail.pointing == RIGHT {
		x = s.tail.x
		y = s.tail.y - 1
	} else if s.tail.pointing == DOWN {
		x = s.tail.x - 1
		y = s.tail.y
	} else {
		x = s.tail.x
		y = s.tail.y + 1
	}
	s.tail.next = &Node{
		x:         x,
		y:         y,
		pointing:  s.tail.pointing,
		prev:      s.tail,
		next:      nil,
		validated: false,
	}
	s.tail = s.tail.next
}
