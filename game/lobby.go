package gsnake

import (
	"sync"
	"turutupa/gsnake/log"
)

var lobby *Lobby
var lock sync.Mutex

type Lobby struct {
	rooms []*Room
}

func NewLobby() *Lobby {
	lock.Lock()
	defer lock.Unlock()

	if lobby != nil {
		return lobby
	}
	lobby = &Lobby{
		rooms: []*Room{},
	}
	return lobby
}

func (l *Lobby) Join(player *Player) *Room {
	lock.Lock()
	defer lock.Unlock()

	var room *Room
	if len(l.rooms) == 0 || l.rooms[len(l.rooms)-1].started || len(l.rooms[len(l.rooms)-1].players) >= MAX_ROOM_SIZE {
		board := NewBoard(ROWS_MULTI, COLS_MULTI)
		game := NewMultiGame(board)
		room = NewRoom(game)
		l.rooms = append(l.rooms, room)
		go room.Run()
		log.Info("No rooms available - created new room: %s", room.id)
	} else {
		room = l.rooms[len(l.rooms)-1]
	}
	room.AddPlayer(player)
	log.Info("Added player %s to room %s", player.name, room.id)
	return room
}
