package gsnake

type Lobby struct {
}

type Room struct {
	players []*Player
	state   string
	game    *Game
}

func NewLobby() *Lobby {
	return &Lobby{}
}
