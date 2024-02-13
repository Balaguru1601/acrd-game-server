package initializers

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func LoadRedis() {

	redisUrl := os.Getenv("REDIS_URL")
	url := redisUrl

	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	RedisClient = redis.NewClient(opts)
}
