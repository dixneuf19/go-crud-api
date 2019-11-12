package greetings

import "github.com/go-redis/redis/v7"

// Redis is a backend greetings storage using redis
type Redis struct {
	client *redis.Client
}

// NewGreetingsRedis create a redis connections
func NewGreetingsRedis(addr string) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Redis{client: client}
}

// Get a key
func (r *Redis) Get(key string) (string, bool) {
	val, err := r.client.Get(key).Result()
	if err == redis.Nil {
		return val, false
	} else if err != nil {
		panic(err)
	}

	if len(val) == 0 {
		return "", false
	}

	return val, true
}

// Add a new greeting
func (r *Redis) Add(key, value string) error {
	err := r.client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// Delete a greeting
func (r *Redis) Delete(key string) error {
	err := r.client.Set(key, "", 0).Err()
	if err != nil {
		return err
	}
	return nil
}
