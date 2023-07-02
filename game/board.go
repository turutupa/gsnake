package gsnake

import (
	"strconv"
	"strings"
)

type Board struct {
	rows   int
	cols   int
	matrix [][]rune
}

func NewBoard(rows, cols int) *Board {
	matrix := Board{}.init(rows, cols)
	board := &Board{rows, cols, matrix}
	board.UpdateLeaderboard(0) // init board to 0
	return board
}

func (b Board) init(rows, cols int) [][]rune {
	var matrix [][]rune
	for i := 0; i < rows; i++ {
		row := make([]rune, cols)
		matrix = append(matrix, row)
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			var cell rune
			switch {
			case i == 0 && j == 0:
				cell = TOP_LEFT
			case i == 0 && j == cols-1:
				cell = TOP_RIGHT
			case i == rows-1 && j == 0:
				cell = BOTTOM_LEFT
			case i == rows-1 && j == cols-1:
				cell = BOTTOM_RIGHT
			case i == 0 || i == rows-1:
				cell = HORIZONTAL
			case j == 0 || j == cols-1:
				cell = VERTICAL
			default:
				cell = ' '
			}
			matrix[i][j] = cell
		}
	}
	return matrix
}

func (b *Board) Restart() {
	newBoard := NewBoard(b.rows, b.cols)
	b.matrix = newBoard.matrix
}

func (b *Board) UpdateFruit(fruit *Fruit) {
	b.matrix[fruit.x][fruit.y] = rune('@')
}

func (b *Board) UpdateLeaderboard(score int) {
	padded_score := strconv.Itoa(score)
	padded_score = strings.Repeat("0", 5-len(padded_score)) + padded_score
	scoreboard := "[ SCORE ]──[ " + padded_score + " ]"
	quitMsg := "[ 'q' to Quit ]"
	padding := 4
	scoreboard = scoreboard + strings.Repeat("─", b.cols-len(scoreboard)-len(quitMsg)-padding) + quitMsg
	j := padding
	for _, r := range scoreboard {
		b.matrix[0][j] = r
		j++
	}
}

func (b *Board) UpdateSnake(node *Node) {
	// render snake
	b.matrix[node.x][node.y] = rune(node.pointing)
	node.render = rune(node.pointing)
	node = node.next

	for node != nil && node.validated {
		var cell rune

		switch {
		case node.next == nil:
			cell = selectSnakeNodeRenderByOrientation(node.pointing)
		case node.pointing == node.next.pointing:
			cell = selectSnakeNodeRenderByOrientation(node.pointing)
		default:
			cell = selectSnakeNodeRenderOnDirection(node)
		}

		b.matrix[node.x][node.y] = cell

		x := node.x
		y := node.y

		if x < 1 || x >= b.rows-1 || y < 1 || y >= b.cols-1 {
			continue
		}

		b.matrix[x][y] = cell
		node.render = cell
		node = node.next
	}
}

func selectSnakeNodeRenderByOrientation(pointing Pointing) rune {
	switch pointing {
	case UP, DOWN:
		return VERTICAL
	default:
		return HORIZONTAL
	}
}

func selectSnakeNodeRenderOnDirection(node *Node) rune {
	switch node.pointing {
	case UP:
		if node.next.pointing == LEFT {
			return BOTTOM_LEFT
		}
		return BOTTOM_RIGHT
	case DOWN:
		if node.next.pointing == LEFT {
			return TOP_LEFT
		}
		return TOP_RIGHT
	case LEFT:
		if node.next.pointing == UP {
			return TOP_RIGHT
		}
		return BOTTOM_RIGHT
	case RIGHT:
		if node.next.pointing == UP {
			return TOP_LEFT
		}
		return BOTTOM_LEFT
	default:
		return 0
	}
}

func (b *Board) Intersects(snake *Snake) bool {
	head := snake.head
	x := head.x
	y := head.y
	if x == 0 || x == b.rows-1 || y == 0 || y == b.cols-1 {
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
