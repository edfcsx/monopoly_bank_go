package accounts

import (
	"github.com/google/uuid"
	"sync"
)

type Account struct {
	Name               string
	Password           string
	Balance            uint64
	AuthenticationHash string
}

var accounts = make(map[string]*Account)
var accountsMutex = sync.Mutex{}

func All() []string {
	accountsMutex.Lock()
	defer accountsMutex.Unlock()

	n := make([]string, 0, len(accounts))

	for _, acc := range accounts {
		n = append(n, acc.Name)
	}

	return n
}

func ExistsByName(name string) *Account {
	accountsMutex.Lock()
	defer accountsMutex.Unlock()

	for _, acc := range accounts {
		if acc.Name == name {
			return acc
		}
	}

	return nil
}

func ExistsByHash(authenticationHash string) *Account {
	accountsMutex.Lock()
	defer accountsMutex.Unlock()

	if acc, ok := accounts[authenticationHash]; ok {
		return acc
	}

	return nil
}

func Create(name string, password string) string {
	accountsMutex.Lock()
	defer accountsMutex.Unlock()

	authenticateHash := uuid.New().String()

	accounts[authenticateHash] = &Account{
		Name:               name,
		Password:           password,
		Balance:            100000,
		AuthenticationHash: authenticateHash,
	}

	return authenticateHash
}
