package gsnake

import (
	"sync"
)

const MAX_ROOM_SIZE int = 6

var lobby *Lobby
var lock *sync.Mutex

type Room struct {
	id              string
	board           *Board
	screen          *Screen
	game            *MultiGame
	lock            *sync.Mutex
	notifyGameStart chan bool
}

type Lobby struct {
	rooms []*Room
	lock  *sync.Mutex
}

func NewLobby() *Lobby {
	lock.Lock()
	defer lock.Unlock()

	if lobby != nil {
		return lobby
	}
	lobby = &Lobby{
		rooms: []*Room{},
		lock:  &sync.Mutex{},
	}
	return lobby
}

func NewRoom(id string, board *Board, screen *Screen, game *MultiGame) *Room {
	return &Room{
		id:              id,
		screen:          screen,
		game:            game,
		lock:            &sync.Mutex{},
		notifyGameStart: make(chan bool),
	}
}

func (l *Lobby) AddRoom(room *Room) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.rooms = append(l.rooms, room)
}

func (r *Room) AddPlayer(player *Player) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.game.players) >= MAX_ROOM_SIZE {
		return false
	}
	r.game.players = append(r.game.players, player)
	r.game.snakes = append(r.game.snakes, NewSnake(r.board))
	if len(r.game.players) == MAX_ROOM_SIZE { // n should be the maximum number of players allowed in a room
		r.notifyGameStart <- true
	}
	return true
}

func (r *Room) Run() {
	for {
		r.game.Setup()
		<-r.notifyGameStart
		r.game.Countdown()
		for {
			r.game.Run()
			if r.game.HasWinner() {
				r.game.Leaderboard()
				break
			}
			r.game.InterRound()
		}
	}
}
