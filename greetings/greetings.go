package greetings

import (
	"fmt"
	"sync"
)

// Provider is a service which store and give greetings for different languages
type Provider interface {
	Get(string) (string, bool)
	Add(string, string) error
	Delete(string) error
}

// Map is a greeting map for a each languages. It is safe for concurrent use with the Lock
type Map struct {
	lock sync.RWMutex
	m    map[string]string
}

// NewGreetingsMap returns a new greetings struct
func NewGreetingsMap() *Map {
	m := &Map{m: make(map[string]string)}
	return m
}

// Get the adequate greeting for the specified language. If there is none, the check flag is false
func (g *Map) Get(lang string) (string, bool) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	greet, ok := g.m[lang]
	return greet, ok
}

// Add or modify a greeting
// return true if succesful
func (g *Map) Add(lang, greet string) error {
	if len(lang) == 0 {
		return fmt.Errorf("empty string is not a valid language")
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	g.m[lang] = greet
	return nil
}

// Delete a entry from Map
func (g *Map) Delete(lang string) error {
	g.lock.Lock()
	defer g.lock.Unlock()
	delete(g.m, lang)
	return nil
}
