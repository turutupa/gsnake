package gsnake

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const VERTICAL rune = '│'
const TOP_LEFT rune = '╭'
const TOP_RIGHT rune = '╮'
const BOTTOM_LEFT rune = '╰'
const BOTTOM_RIGHT rune = '╯'
const HORIZONTAL rune = '─'

const MAX_HIGH_SCORE_DIGITS = 4

type Screen struct {
	writer        io.Writer
	isFirstRender bool
	termRows      int
	termCols      int
}

func NewScreen(writer io.Writer) *Screen {
	return &Screen{writer, true, 30, 50}
}

func (s *Screen) Restart() {
	s.isFirstRender = true
}

func (s *Screen) Clear() {
	s.writer.Write([]byte("\033[H\033[2J"))
}

func (s *Screen) SetSize(rows int, cols int) {
	s.termRows = rows
	s.termCols = cols
}

func (s *Screen) InputBox(label string, value string) {
	row := s.termRows / 3
	startLine := s.termCols/2 - len(label)/2
	s.printBold(row, startLine, label)
	row++
	startLine = s.termCols/2 - MAX_PLAYER_NAME_LEN/2
	if len(value) < MAX_PLAYER_NAME_LEN {
		value = value + "_" + strings.Repeat(" ", MAX_PLAYER_NAME_LEN-len(value))
		s.print(row, startLine, value)
	} else {
		s.print(row, startLine, value)
	}
	row = row + 2
	continueLabel := "Press enter to continue"
	startLine = s.termCols/2 - len(continueLabel)/2
	s.print(row, startLine, continueLabel)
}

func (s *Screen) RenderMenu(termSize TermSize, title string, options []string, selected int) {
	startLine := termSize.cols/2 - len(title)/2 - 1
	row := termSize.rows / 3
	s.printBold(row, startLine, title)
	row += 2
	optionIndex := 0
	for _, game := range options {
		paddingRight := 8
		if optionIndex == selected {
			selectedIndicatorLeft := "* "
			selectedIndicatorRight := " *"
			gameFmt := strings.Repeat(" ", len(title)-len(selectedIndicatorLeft)-(len(game)/2)-paddingRight)
			gameFmt = gameFmt + selectedIndicatorLeft + game + selectedIndicatorRight
			gameFmt = gameFmt + strings.Repeat(" ", paddingRight-len(selectedIndicatorRight))
			s.printBold(row, startLine, gameFmt)
		} else {
			gameFmt := strings.Repeat(" ", len(title)-(len(game)/2)-paddingRight)
			gameFmt = gameFmt + game
			gameFmt = gameFmt + strings.Repeat(" ", paddingRight)
			for i, r := range gameFmt {
				s.printChar(row, startLine+i, r)
			}
		}
		if game == INSANITY {
			row += 2
		} else {
			row++
		}
		optionIndex++
	}
}

func (s *Screen) RenderBoard(board *Board) {
	offsetRow, offsetCol := s.offsets(board)
	for i := 0; i < board.rows; i++ {
		for j := 0; j < board.cols; j++ {
			s.printChar(offsetRow+i, offsetCol+j, board.matrix[i][j])
		}
	}
}

func (s *Screen) RenderBoardFrame(board *Board) {
	offsetRow, offsetCol := s.offsets(board)
	for i := 0; i < board.cols; i++ {
		s.printChar(offsetRow+0, offsetCol+i, board.matrix[0][i])
		s.printChar(offsetRow+board.rows-1, offsetCol+i, board.matrix[board.rows-1][i])
	}
	for i := 0; i < board.rows; i++ {
		s.printChar(offsetRow+i, offsetCol+0, board.matrix[i][0])
		s.printChar(offsetRow+i, offsetCol+board.cols-1, board.matrix[i][board.cols-1])
	}
}

func (s *Screen) RenderCountdown(board *Board, label string, n int) {
	label = fmt.Sprintf("%s %d", label, n)
	offsetRow, offsetCol := s.offsets(board)
	row := board.rows / 2
	col := board.cols/2 - len(label)/2
	s.printBold(offsetRow+row, offsetCol+col, label)
}

func (s *Screen) RenderWarning(board *Board, txt string) {
	ro, co := s.offsets(board)
	s.printBold(ro+board.rows/2-1, co+board.cols/2-len(txt)/2, txt)
}

func (s *Screen) RenderFruit(board *Board, fruit *Fruit) {
	offsetRow, offsetCol := s.offsets(board)
	s.printChar(offsetRow+fruit.x, offsetCol+fruit.y, '@')
}

// used in solo game
func (s *Screen) RenderScore(board *Board, score int) {
	// - top border -
	// we only want to re-render the score so
	// calcualte where the score is positioned,
	// which is between the second set of brackets
	offsetRow, offsetCol := s.offsets(board)
	bracket_counter := 0
	for j := 0; j < board.cols; j++ {
		if board.matrix[0][j] == '[' || board.matrix[0][j] == ']' {
			bracket_counter = bracket_counter + 1
			continue
		}
		if bracket_counter == 3 {
			s.printChar(offsetRow, offsetCol+j, board.matrix[0][j])
		}
		if bracket_counter >= 4 {
			break
		}
	}
}

func (s *Screen) Remove(board *Board, head *Node) {
	s.printChar(s.offsetRow(board)+head.row, s.offsetCol(board)+head.col, ' ')
	last := head
	for head != nil && head.validated {
		last = head
		head = head.next
	}
	s.printChar(s.offsetRow(board)+last.row, s.offsetCol(board)+last.col, ' ')
}

// room for optimization
func (s *Screen) RenderSnake(board *Board, snake *Snake) {
	head := snake.head
	offsetRow, offsetCol := s.offsets(board)
	s.print(
		offsetRow+head.row,
		offsetCol+head.col,
		color(string(head.render), snake.color),
	)
	node := head.next
	for node != nil && node.validated {
		s.print(
			offsetRow+node.row,
			offsetCol+node.col,
			color(string(node.render), snake.color),
		)
		node = node.next
	}
}

func (s *Screen) RenderSingleLeaderboard(board *Board, difficulty string, scores []*Player, player *Player) {
	spacing := strings.Repeat(" ", 4)
	title := "|" + spacing + "TOP SCORES" + spacing + "|"
	// how many high scores to render, if there is a new one,
	// that one is going to be the fifth top score
	var leaderboarSize int
	if player != nil {
		leaderboarSize = 4
	} else {
		leaderboarSize = 5
	}
	for len(scores) < leaderboarSize {
		// this is just to make sure we always render 5 rows of scores,
		// even if player hasn't already played 5 times
		scores = append(scores, NewPlayer(""))
	}
	scores = scores[:leaderboarSize]
	marginLeft := len(title)/2 - 1
	startLine := board.cols/2 - marginLeft
	row := board.rows/3 + 1
	offsetRow, offsetCol := s.offsets(board)

	// print top border
	for i := 0; i < len(title); i++ {
		position := startLine + i
		switch i {
		case 0:
			s.printChar(offsetRow+row, offsetCol+position, TOP_LEFT)
		case len(title) - 1:
			s.printChar(offsetRow+row, offsetCol+position, TOP_RIGHT)
		default:
			s.printChar(offsetRow+row, offsetCol+position, HORIZONTAL)

		}
	}
	row++

	// print title
	for i, r := range title {
		s.printChar(offsetRow+row, offsetCol+startLine+i, r)
	}
	row++

	// print difficulty
	padding := strings.Repeat(" ", (len(title)-len(difficulty)-2)/2)
	diff := "|" + padding + difficulty + padding + "|"
	for i, r := range diff {
		s.printChar(offsetRow+row, offsetCol+startLine+i, r)
	}
	row++

	// print empty line
	emptyLine := "|" + strings.Repeat(" ", len(title)-2) + "|"
	for i, r := range emptyLine {
		s.printChar(offsetRow+row, offsetCol+startLine+i, r)
	}
	row++

	// print scores
	renderedNewHighScore := false
	i := 0
	for i < len(scores) {
		score := scores[i]
		var scoreFmt string
		var sc int
		var pl string
		isHighScore := player != nil && !renderedNewHighScore && player.score > score.score
		if isHighScore && !renderedNewHighScore {
			sc = player.score
			pl = player.name
			renderedNewHighScore = true
		} else {
			sc = score.score
			pl = score.name
			i++
		}

		rightPadding := 2
		leftPadding := 2
		scoreStr := strconv.Itoa(sc)
		scoreFmt = strings.Repeat(" ", leftPadding) + pl
		if len(pl) < MAX_PLAYER_NAME_LEN && isHighScore {
			// format when submitting high score
			cursor := "_"
			scoreFmt = scoreFmt + cursor + strings.Repeat(" ", MAX_PLAYER_NAME_LEN-len(pl)-1)
		} else {
			// add padding spaces for player names
			scoreFmt = scoreFmt + strings.Repeat(" ", MAX_PLAYER_NAME_LEN-len(pl))
		}
		// add space between names and actual scores
		scoreFmt = scoreFmt + strings.Repeat(" ", len(title)-len(scoreFmt)-MAX_HIGH_SCORE_DIGITS-rightPadding)
		// add leading zeroes to score
		scoreFmt = scoreFmt + strings.Repeat("0", MAX_HIGH_SCORE_DIGITS-len(strconv.Itoa(sc)))
		scoreFmt = scoreFmt + scoreStr
		scoreFmt = scoreFmt + strings.Repeat(" ", rightPadding) // add padding to the right if needed

		for j, r := range scoreFmt {
			position := startLine + j
			if j == 0 || j == len(title)-1 {
				s.printChar(offsetRow+row, offsetCol+position, '|')
			} else {
				s.printChar(offsetRow+row, offsetCol+position, r)
			}
		}
		row++
	}

	// render border bottom
	for i := range title {
		position := startLine + i
		if i == 0 {
			s.printChar(offsetRow+row, offsetCol+position, BOTTOM_LEFT)
		} else if i == len(title)-1 {
			s.printChar(offsetRow+row, offsetCol+position, BOTTOM_RIGHT)
		} else {
			s.printChar(offsetRow+row, offsetCol+position, HORIZONTAL)
		}
	}
	row = row + 2
	var msg string
	if player != nil {
		msg = " PRESS ENTER TO SUBMIT NEW HIGH SCORE "
	} else {
		msg = " PRESS ENTER TO CONTINUE "
	}
	s.printBold(offsetRow+row, offsetCol+board.cols/2-len(msg)/2, msg)
}

func (s *Screen) RenderMultiLeaderboard(board *Board, players []*Player, hasWinner bool) {
	offsetRow, offsetCol := s.offsets(board)
	var title string
	if hasWinner {
		title = fmt.Sprintf("Congratulations %s!", players[0].name)
	} else {
		title = "~~ Leaderboard ~~"
	}
	row := offsetRow + board.rows/2
	col := offsetCol + board.cols/2 - len(title)/2
	s.printBold(row, col, title)
	row++
	for i, player := range players {
		score := strings.Repeat("0", MAX_HIGH_SCORE_DIGITS-len(strconv.Itoa(player.score))) + strconv.Itoa(player.score)
		name := player.name
		separator := strings.Repeat(" ", len(title)-3-len(name)-len(score))
		txt := fmt.Sprintf("%d. %s%s%s", i, name, separator, score)
		s.printBold(row, col, txt)
		row++
	}
}

func (s *Screen) HideCursor() {
	fmt.Fprint(s.writer, "\033[?25l")
}

func (s *Screen) ShowCursor() {
	fmt.Fprint(s.writer, "\033[?25h")
}

func (s *Screen) GameOver() {
	s.writer.Write([]byte(GAME_OVER))
}

func (s *Screen) PrintLogo() {
	s.writer.Write([]byte(LOGO))
}

func color(str, color string) string {
	return color + str + reset
}

func (s *Screen) print(row, col int, str string) {
	fmt.Fprintf(s.writer, "\033[%d;%dH%s", row+1, col+1, str)
}

func (s *Screen) printBold(row, col int, r string) {
	fmt.Fprintf(s.writer, "\033[%d;%dH\033[1m%s\033[0m", row+1, col+1, r)
}

func (s *Screen) printChar(row, col int, r rune) {
	fmt.Fprintf(s.writer, "\033[%d;%dH%c", row+1, col+1, r)
}

func (s *Screen) offsets(board *Board) (int, int) {
	return s.offsetRow(board), s.offsetCol(board)
}

func (s *Screen) offsetRow(board *Board) int {
	return (s.termRows - board.rows) / 2
}

func (s *Screen) offsetCol(board *Board) int {
	return (s.termCols - board.cols) / 2
}
