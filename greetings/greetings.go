package greetings

var greetings = make(map[string]string)

// GetGreetings returns the greetings map
func GetGreetings() map[string]string {
	return greetings
}

// AddGreeting add or modify a greeting
// return true if succesful
func AddGreeting(lang, greet string) bool {
	if len(lang) == 2 {
		greetings[lang] = greet
		return true
	}
	return false
}
