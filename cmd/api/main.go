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

	// This is obviousy not production code :)
	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "https://petstore.swagger.io")
			header.Set("Access-Control-Allow-Headers", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})
	router.Handler(http.MethodPost, "/users", corsMiddleware(auth.CreateUserHandler(userRepo)))
	router.Handler(http.MethodGet, "/players", corsMiddleware(authMiddleware(http.HandlerFunc(playerHandle.GetAllPlayers))))
	router.Handler(http.MethodGet, "/players/:id", corsMiddleware(authMiddleware(http.HandlerFunc(playerHandle.GetPlayerByID))))
	router.Handler(http.MethodPost, "/players", corsMiddleware(authMiddleware(http.HandlerFunc(playerHandle.CreatePlayer))))
	router.Handler(http.MethodPost, "/players/:id", corsMiddleware(authMiddleware(http.HandlerFunc(playerHandle.UpdatePlayer))))
	router.Handler(http.MethodDelete, "/players/:id", corsMiddleware(authMiddleware(http.HandlerFunc(playerHandle.DeletePlayer))))

	log.Fatal(http.ListenAndServe(":"+config.HttpPort, router))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
		header.Set("Access-Control-Allow-Origin", "https://petstore.swagger.io")
		header.Set("Access-Control-Allow-Headers", "*")

		next.ServeHTTP(w, r)
	})
}
