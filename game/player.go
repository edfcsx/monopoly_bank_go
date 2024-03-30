package Game

import "sync"

type Player struct {
	Name     string
	Password string
	Balance  int
}

var Players = make(map[string]*Player)
var PlayersMutex = sync.Mutex{}

func PlayerExists(name string) *Player {
	PlayersMutex.Lock()
	defer PlayersMutex.Unlock()

	if player, ok := Players[name]; ok {
		return player
	}

	return nil
}

func CreatePlayer(name, password string) *Player {
	PlayersMutex.Lock()
	defer PlayersMutex.Unlock()

	player := &Player{
		Name:     name,
		Password: password,
		Balance:  100000,
	}

	Players[name] = player
	return player
}
