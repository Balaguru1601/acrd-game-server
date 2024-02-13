package main

// import "fmt"
import (
	"fmt"
	"go-backend/initializers"
	"go-backend/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadRedis()
}

type Member struct {
	Username string `json:"username"`
	Score    string `json:"score"`
}

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://card-game-puce.vercel.app"}

	r.Use(cors.New(config))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello there",
		})
	})

	r.GET("/leaderboard", func(c *gin.Context) {
		userData := redis.GetAllValues(c)
		c.JSON(200, gin.H{
			"message": "Success",
			"data":    userData,
		})
	})

	r.POST("/register", func(c *gin.Context) {

		var userData struct {
			Username string
			Secret   string
		}

		if err := c.Bind(&userData); err != nil {
			fmt.Println(err)
		}

		userExists, e := redis.CheckUserExists(c, userData.Username)
		if e != nil || userExists {
			fmt.Println(e, userExists)
			c.JSON(200, gin.H{
				"message": "Username already exists!",
				"success": false,
			})
			return
		}

		err := redis.SetSecretValue(c, userData.Username, userData.Secret)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Something went wrong",
				"success": false,
			})
			panic(err)
		}

		redis.SetValue(c, userData.Username, "0")

		c.JSON(200, gin.H{
			"message": "Success",
			"success": true,
			"score":   "0",
		})
	})

	r.POST("/verify-user", func(c *gin.Context) {
		var userData struct {
			Username string
			Secret   string
		}
		c.Bind(&userData)
		fmt.Println(userData)
		validity, err := redis.CheckSecretValue(c, userData.Username, userData.Secret)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Something went wrong",
				"success": false,
			})
			panic(err)
		}

		if validity {

			score, err := redis.GetValue(c, userData.Username)
			if err != nil || score == "" {
				c.JSON(200, gin.H{
					"message": "Success",
					"score":   "0",
					"success": true,
				})
			} else {

				c.JSON(200, gin.H{
					"message": "Success",
					"success": true,
					"score":   score,
				})
			}
			return
		} else {
			c.JSON(400, gin.H{
				"message": "Unauthorized",
				"success": false,
			})
			return
		}

	})

	r.POST("/set-score", func(c *gin.Context) {
		var userData struct {
			Username string
			Score    string
		}
		c.Bind(&userData)

		err := redis.SetValue(c, userData.Username, userData.Score)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Something went wrong",
				"success": false,
			})
			panic(err)
		}

		c.JSON(200, gin.H{
			"message":  "Success",
			"success":  true,
			"username": userData.Username,
			"score":    userData.Score,
		})

	})

	r.POST("/save-game", func(c *gin.Context) {

		var gameData redis.GameData
		if err := c.ShouldBindJSON(&gameData); err != nil {
			c.JSON(200, gin.H{"message": "Something went wrong!", "success": false})
			return
		}
		res := redis.SetGameData(c, gameData)
		if !res {
			c.JSON(500, gin.H{"message": "Error storing data", "success": false})
			return
		}

		c.JSON(200, gin.H{"message": "Data stored successfully", "success": true})

	})

	r.GET("/get-game/:username", func(c *gin.Context) {
		username := c.Param("username")

		data, err := redis.GetGameData(c, username)
		if err != nil {
			switch err {
			case redis.RedisNil:
				c.JSON(200, gin.H{"message": "Username not found", "success": false})
			default:
				c.JSON(200, gin.H{"message": "Error retrieving data", "success": false})
			}
			return
		}

		c.JSON(200, gin.H{"message": "Data retrieved!", "success": true, "data": data})
	})
	r.Run()
}
