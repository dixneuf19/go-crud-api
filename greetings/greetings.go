package greetings

// Provider is a service which store and give greetings for different languages
type Provider interface {
	Get(string) (string, bool)
	Add(string, string) error
	Delete(string) error
}
