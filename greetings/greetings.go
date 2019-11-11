package greetings

import (
	"fmt"
)

// Provider is a service which store and give greetings for different languages
type Provider interface {
	Get(string) (string, bool)
	Add(string, string) error
	Delete(string) error
}

// GreetingsMap is a greeting for a specified language
type GreetingsMap map[string]string

// NewGreetingsMap returns a new greetings struct
func NewGreetingsMap() GreetingsMap {
	return make(map[string]string)
}

// Get the adequate greeting for the specified language. If there is none, the check flag is false
func (g GreetingsMap) Get(lang string) (string, bool) {
	greet, ok := g[lang]
	return greet, ok
}

// Add or modify a greeting
// return true if succesful
func (g GreetingsMap) Add(lang, greet string) error {
	if len(lang) == 0 {
		return fmt.Errorf("empty string is not a valid language")
	}
	g[lang] = greet
	return nil
}

// Delete a entry from greetingsMap
func (g GreetingsMap) Delete(lang string) error {
	delete(g, lang)
	return nil
}
