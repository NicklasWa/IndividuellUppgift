package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"net/url"

	"github.com/gin-gonic/gin"
	"systementor.se/cloudgolangapi/data"
)

type PageView1 struct {
	Title  string
	Rubrik string
}

type PageView struct {
	CurrentUser string
	PageTitle   string
	Title       string
	Text        string
}

type LoginView struct {
	CurrentUser string
	PageTitle   string
	Error       bool
	Email       string
}

type GameHomeView struct {
	CurrentUser string
	PageTitle   string
	Error       bool
	Email       string
	History     string
}

type GameView struct {
	CurrentUser string
	PageTitle   string
	Error       bool
	Email       string
	Result      string
	History     string
}

var gameState int = 0

var config Config
var theRandom *rand.Rand

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "gameHome.html", &GameView{PageTitle: "Spela"})
	gameState = 0
}

func homePost(c *gin.Context) {
	var winner string = ""
	switch {
	case gameState == 0:
		c.Status(200)
		c.HTML(http.StatusOK, "gamePlay.html", &GameView{PageTitle: "Välj ett alternativ"})
		gameState = 1
		break

	case gameState == 1:
		button := c.PostForm("button")
		switch {
		case button == "stenButton":
			winner = play("STONE")
			break
		case button == "saxButton":
			winner = play("SCISSOR")
			break
		case button == "paseButton":
			winner = play("BAG")
			break
		}

		totalGames, wins := data.Stats()

		c.Status(200)
		var hist string = fmt.Sprintf("Totalt antal spel: %d. Antal vunna spel: %d.", totalGames, wins)
		c.HTML(http.StatusOK, "gameDone.html", &GameView{PageTitle: "Resultat", Result: "Segrare är " + winner, History: hist})
		gameState = 2
		break

	case gameState == 2:
		c.Status(200)
		c.HTML(http.StatusOK, "gamePlay.html", &GameView{PageTitle: "Välj ett alternativ"})
		gameState = 1
		break
	}
}

func gameHome(c *gin.Context) {
	//c.HTML(http.StatusOK, "gameHome.html", &GameHomeView{PageTitle: "Game", History: "Dator - Spelare.."})
	c.HTML(http.StatusOK, "gameHome.html", &GameView{PageTitle: "Spela", Result: "Resultat:"})
}

func gameHomePost(c *gin.Context) {
	if c.PostForm("button") == "startGameButton" {
		//c.Status(200)
		//c.HTML(http.StatusOK, "game.html", &LoginView{PageTitle: "Börja!"})

		location := url.URL{Path: "/game"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}
}

func game(c *gin.Context) {
	c.HTML(http.StatusOK, "game.html", &LoginView{PageTitle: "Game"})
}

func gamePost(c *gin.Context) {
	button := c.PostForm("button")
	winner := ""
	switch {
	case button == "stenButton":
		winner = play("STONE")
		break
	case button == "staxButton":
		winner = play("SCISSOR")
		break
	case button == "paseButton":
		winner = play("BAG")
		break
	}

	c.Status(200)
	c.HTML(http.StatusOK, "gameHome.html", &GameView{PageTitle: "Spela igen", Result: "Resultat: " + winner})

}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", &LoginView{PageTitle: "Login"})
}

func enableCors(c *gin.Context) {
	(*c).Header("Access-Control-Allow-Origin", "*")
}

func apiStats(c *gin.Context) {
	enableCors(c)
	totalGames, wins := data.Stats()
	c.JSON(http.StatusOK, gin.H{"totalGames": totalGames, "wins": wins})
}

func play(yourSelection string) string {
	mySelection := randomizeSelection()
	winner := "Oavgjort"
	if yourSelection == "STONE" && mySelection == "SCISSOR" {
		winner = "Du"
	}
	if yourSelection == "SCISSOR" && mySelection == "BAG" {
		winner = "Du"
	}
	if yourSelection == "BAG" && mySelection == "STONE" {
		winner = "Du"
	}
	if mySelection == "STONE" && yourSelection == "SCISSOR" {
		winner = "Dator"
	}
	if mySelection == "SCISSOR" && yourSelection == "BAG" {
		winner = "Dator"
	}
	if mySelection == "BAG" && yourSelection == "STONE" {
		winner = "Dator"
	}
	data.SaveGame(yourSelection, mySelection, winner)
	return winner
}

func apiPlay(c *gin.Context) {
	enableCors(c)
	yourSelection := c.Query("yourSelection")
	mySelection := randomizeSelection()
	winner := "Tie"
	if yourSelection == "STONE" && mySelection == "SCISSOR" {
		winner = "You"
	}
	if yourSelection == "SCISSOR" && mySelection == "BAG" {
		winner = "You"
	}
	if yourSelection == "BAG" && mySelection == "STONE" {
		winner = "You"
	}
	if mySelection == "STONE" && yourSelection == "SCISSOR" {
		winner = "Computer"
	}
	if mySelection == "SCISSOR" && yourSelection == "BAG" {
		winner = "Computer"
	}
	if mySelection == "BAG" && yourSelection == "STONE" {
		winner = "Computer"
	}
	data.SaveGame(yourSelection, mySelection, winner)
	c.JSON(http.StatusOK, gin.H{"winner": winner, "yourSelection": yourSelection, "computerSelection": mySelection})
}

func randomizeSelection() string {
	val := theRandom.Intn(3) + 1
	if val == 1 {
		return "STONE"
	}
	if val == 2 {
		return "SCISSOR"
	}
	if val == 3 {
		return "BAG"
	}
	return "ERROR"

}

func main() {
	theRandom = rand.New(rand.NewSource(time.Now().UnixNano()))
	readConfig(&config)

	data.InitDatabase(config.Database.File,
		config.Database.Server,
		config.Database.Database,
		config.Database.Username,
		config.Database.Password,
		config.Database.Port)

	router := gin.Default()
	router.LoadHTMLGlob("templates/**")
	router.GET("/", home)
	router.POST("/", homePost)
	//router.GET("/game/", game)
	//router.POST("/game/", gamePost)
	//router.GET("/login", login)
	//router.GET("/api/play", apiPlay)
	//router.GET("/api/stats", apiStats)

	// router.GET("/api/employee/:id", apiEmployeeById)
	// router.PUT("/api/employee/:id", apiEmployeeUpdateById)
	// router.DELETE("/api/employee/:id", apiEmployeeDeleteById)
	// router.POST("/api/employee", apiEmployeeAdd)

	// router.GET("/api/employees", employeesJson)
	// router.GET("/api/addemployee", addEmployee)
	// router.GET("/api/addmanyemployees", addManyEmployees)
	router.Run(":8080")

}
