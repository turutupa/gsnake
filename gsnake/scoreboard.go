package gsnake

type Scoreboard struct {
	storageDir string
	scoreboard []int
}

func NewScoreboard(storageDir string) *Scoreboard {
	scoreboard := &Scoreboard{storageDir, []int{}}
	scoreboard.read()
	return scoreboard
}

func (s *Scoreboard) read() {}

func (s *Scoreboard) get() []int {
	return []int{}
}
