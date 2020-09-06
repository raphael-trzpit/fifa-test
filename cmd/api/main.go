package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	"github.com/raphael-trzpit/fifa-test/internal/auth"
	"github.com/raphael-trzpit/fifa-test/internal/players"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err.Error())
	}

	db, err := gorm.Open(mysql.Open(config.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepo, err := auth.NewUserRepositoryMysl(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	playerRepo, err := players.NewPlayerRepositoryMysl(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	playerHandle := players.Handle{Repository: playerRepo}
	authMiddleware := auth.AuthMiddleware(userRepo)

	router := httprouter.New()
	router.Handler(http.MethodPost, "/users", auth.CreateUserHandler(userRepo))
	router.Handler(http.MethodGet, "/players", authMiddleware(http.HandlerFunc(playerHandle.GetAllPlayers)))
	router.Handler(http.MethodGet, "/players/:id", authMiddleware(http.HandlerFunc(playerHandle.GetPlayerByID)))
	router.Handler(http.MethodPost, "/players", authMiddleware(http.HandlerFunc(playerHandle.CreatePlayer)))
	router.Handler(http.MethodPost, "/players/:id", authMiddleware(http.HandlerFunc(playerHandle.UpdatePlayer)))
	router.Handler(http.MethodDelete, "/players/:id", authMiddleware(http.HandlerFunc(playerHandle.DeletePlayer)))

	log.Fatal(http.ListenAndServe(":"+config.HttpPort, router))
}
