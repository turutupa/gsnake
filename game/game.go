package gsnake

type Game interface {
	Run()
	Leaderboard()
	Restart()
	SetDifficulty(string)
	Stop()
	Strategy(event rune)
}
