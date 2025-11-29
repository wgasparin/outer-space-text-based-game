package main

import (
	"fmt"
	"net/http"
	character "outerspace/character"
	singleton "outerspace/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/someJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}

		// will output : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(http.StatusOK, data)
	})

	router.POST("/players", createPlayer)

	router.GET("/players", func(c *gin.Context) {
		players := singleton.GetInstance[PlayerDataSource]()
		c.JSON(http.StatusOK, gin.H{
			"players": players.Players,
		})

	})

	router.Run() // listens on 0.0.0.0:8080 by default
}

func createPlayer(c *gin.Context) {

	player := Player{}

	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var striker *character.Chararcter

	if player.Type == character.Striker {
		// NewStriker is defined with a receiver on *Chararcter in the
		// character package, so call it on a zero receiver to get a
		// properly initialized Striker instance.
		striker = (&character.Chararcter{}).NewStriker(player.Name, 1)
		// Other types to be implemented
	}

	players := singleton.GetInstance[PlayerDataSource]()
	players.Players = append(players.Players, player)

	var x character.CharacterType = character.Striker

	fmt.Println(x)

	c.JSON(http.StatusCreated, gin.H{
		"status": "player created" + striker.String(),
	})
}

type Player struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

type PlayerDataSource struct {
	Players []Player
}
