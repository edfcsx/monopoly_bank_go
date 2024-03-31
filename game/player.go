package Game

import (
	"github.com/google/uuid"
	"sync"
)

type Player struct {
	Name     string
	Password string
	Balance  int
	AuthHash string
}

var Players = make(map[string]*Player)
var PlayersMutex = sync.Mutex{}

func PlayerExistsByName(name string) *Player {
	PlayersMutex.Lock()
	defer PlayersMutex.Unlock()

	for _, player := range Players {
		if player.Name == name {
			return player
		}
	}

	return nil
}

func PlayerExistsByHash(hash string) *Player {
	PlayersMutex.Lock()
	defer PlayersMutex.Unlock()

	if player, ok := Players[hash]; ok {
		return player
	}

	return nil
}

func CreatePlayer(name, password string) string {
	PlayersMutex.Lock()
	defer PlayersMutex.Unlock()

	hash := uuid.New().String()

	player := &Player{
		Name:     name,
		Password: password,
		Balance:  100000,
		AuthHash: hash,
	}

	Players[hash] = player
	return hash
}
