package handlers

import (
	"strings"
	"sync"
)

var (
	stateMu sync.RWMutex
	state   = make(map[int64]string)
)

func setState(chatID int64, s string) {
	stateMu.Lock()
	defer stateMu.Unlock()
	state[chatID] = s
}

func getState(chatID int64) string {
	stateMu.RLock()
	defer stateMu.RUnlock()
	return state[chatID]
}

func clearState(chatID int64) {
	stateMu.Lock()
	defer stateMu.Unlock()
	delete(state, chatID)
}

func trimLower(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
