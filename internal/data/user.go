package data

import "sync"

type User struct {
	APIKey string
	State  bool
}

// Global map with a mutex for thread safety
var (
	userMap = make(map[string]User)
	mu      sync.RWMutex
)

func SetUser(userID, apiKey string, state bool) {
	mu.Lock()
	defer mu.Unlock()
	userMap[userID] = User{
		APIKey: apiKey,
		State:  state,
	}
}

func GetUser(userID string) (User, bool) {
	mu.RLock()
	defer mu.RUnlock()
	data, ok := userMap[userID]
	return data, ok
}

func UpdateUserState(userID string, state bool) {
	mu.Lock()
	defer mu.Unlock()

	if user, ok := userMap[userID]; ok {
		user.State = state
		userMap[userID] = user
	}
}
