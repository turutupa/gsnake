package gsnake

type Pointing rune

const (
	UP    Pointing = '▲'
	DOWN           = '▼'
	RIGHT          = '►'
	LEFT           = '◄'
)

type Snake struct {
	head *Node
	tail *Node
}

type Node struct {
	x         int
	y         int
	pointing  Pointing
	prev      *Node
	next      *Node
	render    rune
	validated bool
}

type Direction struct {
	X int
	Y int
}

func NewSnake(screen *Screen) *Snake {
	snake := &Snake{}
	snake.head = &Node{
		x:         screen.rows / 2,
		y:         screen.cols / 5,
		pointing:  RIGHT,
		prev:      nil,
		next:      nil,
		render:    HORIZONTAL,
		validated: true,
	}
	node := snake.head
	for i := 0; i < 6; i++ { // initial length of 7
		node.next = &Node{
			x:         node.x,
			y:         node.y - 1,
			pointing:  RIGHT,
			prev:      node,
			next:      nil,
			render:    HORIZONTAL,
			validated: true,
		}
		node = node.next
	}
	snake.tail = node
	return snake
}

func (s *Snake) restart(screen *Screen) {
	newSnake := NewSnake(screen)
	s.head = newSnake.head
	s.tail = newSnake.tail
}

func (s *Snake) move() {
	node := s.head
	x_prev := node.x
	y_prev := node.y
	pointing_prev := node.pointing
	if node.pointing == UP {
		node.x = node.x - 1
	} else if node.pointing == RIGHT {
		node.y = node.y + 1
	} else if node.pointing == LEFT {
		node.y = node.y - 1
	} else {
		node.x = node.x + 1
	}
	node = node.next
	for node != nil {
		pointing_tmp := node.pointing
		x_tmp := node.x
		y_tmp := node.y
		node.x = x_prev
		node.y = y_prev
		node.pointing = pointing_prev
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
