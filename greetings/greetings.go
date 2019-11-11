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

// Add or modify a greeting
// return true if succesful
func (g Greetings) Add(lang, greet string) error {
	if len(lang) == 0 {
		return fmt.Errorf("empty string is not a valid language")
	}
	g[lang] = greet
	return nil
}
