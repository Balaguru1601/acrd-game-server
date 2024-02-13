package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"go-backend/initializers"
	"sort"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var groupName string = "go:Users"

type Member struct {
	Username string `json:"username"`
	Score    string `json:"score"`
}

type GameData struct {
	Username string          `json:"username"`
	Data     json.RawMessage `json:"data"`
}

func (d *GameData) MarshalBinary() ([]byte, error) {
	return json.Marshal(d.Data)
}

var RedisNil = redis.Nil

func SetValue(ctx context.Context, username string, value string) error {
	_, err := initializers.RedisClient.HSet(context.Background(), groupName, username, value).Result()
	if err != nil {
		panic(err)
	}
	if err != nil {
		return err
	}

	return nil
}

func GetValue(ctx context.Context, username string) (string, error) {
	score, err := initializers.RedisClient.HGet(ctx, groupName, username).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	return score, nil
}

func GetAllValues(ctx context.Context) []Member {
	members, err := initializers.RedisClient.HGetAll(context.Background(), groupName).Result()
	if err != nil {
		panic(err)
	}

	memberArray := make([]Member, 0, len(members))

	for username, scoreStr := range members {
		member := Member{
			Username: username,
			Score:    scoreStr,
		}
		memberArray = append(memberArray, member)
	}

	sort.SliceStable(memberArray, func(i, j int) bool {
		score1, err := strconv.Atoi(memberArray[i].Score)
		if err != nil {
			panic(err)
		}
		score2, err := strconv.Atoi(memberArray[j].Score)
		if err != nil {
			panic(err)
		}

		return score1 > score2
	})

	return memberArray

}

var secretGrpName = "go:Secret"

func SetSecretValue(ctx context.Context, username string, value string) error {
	fmt.Println(username, value)
	_, err := initializers.RedisClient.HSet(context.Background(), secretGrpName, username, value).Result()
	if err != nil {
		return err
	}

	return nil
}

func CheckSecretValue(ctx context.Context, username string, givenSecret string) (bool, error) {
	secret, err := initializers.RedisClient.HGet(ctx, secretGrpName, username).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	return secret == givenSecret, nil
}

func CheckUserExists(ctx context.Context, username string) (bool, error) {
	secret, err := initializers.RedisClient.HExists(ctx, secretGrpName, username).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println(err)
			return true, nil
		}
	}
	if secret {
		return true, nil
	} else {
		return false, nil
	}
}

const gameGrp = "go:game"

func SetGameData(ctx context.Context, gameData GameData) bool {
	jsonData, e := json.Marshal(gameData.Data)
	if e != nil {
		return false
	}
	err := initializers.RedisClient.HSet(ctx, gameGrp, gameData.Username, jsonData).Err()
	fmt.Println(err)
	return err == nil
}

func GetGameData(ctx context.Context, username string) (interface{}, error) {
	val, err := initializers.RedisClient.HGet(ctx, gameGrp, username).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(val), &jsonData); err != nil {
		return nil, err
	}

	return jsonData, nil
}
