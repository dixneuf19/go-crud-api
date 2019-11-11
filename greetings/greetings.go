package greetings

import (
	"fmt"
)

// Greetings is a greeting for a specified language
type Greetings map[string]string

// NewGreetings returns a new greetings struct
func NewGreetings() Greetings {
	return make(map[string]string)
}

// Get the adequate greeting for the specified language. If there is none, the check flag is false
func (g Greetings) Get(lang string) (string, bool) {
	greet, ok := g[lang]
	return greet, ok
}

// Add or modify a greeting
// return true if succesful
func (g Greetings) Add(lang, greet string) error {
	if len(lang) == 0 {
		return fmt.Errorf("empty string is not a valid language")
	}
	g[lang] = greet
	return nil
}

// Delete a entry from greetings
func (g Greetings) Delete(lang string) error {
	delete(g, lang)
	return nil
}
