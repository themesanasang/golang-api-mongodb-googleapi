package main

import (
	"log"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	//"github.com/labstack/gommon/log"
	"github.com/themesanasang/testdb/handler"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	e := echo.New()
	//e.Logger.SetLevel(log.ERROR)
	e.Use(middleware.Logger())
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(handler.Key),
		Skipper: func(c echo.Context) bool {
			// Skip authentication for and signup login requests
			if c.Path() == "/login" {
				return true
			}
			return false
		},
	}))

	//CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{"localhost:27017"},
		Timeout:  60 * time.Second,
		Database: "",
		Username: "",
		Password: "",
	}

	// to our MongoDB.
	db, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create indices
	if err = db.Copy().DB("pokemondb").C("user").EnsureIndex(mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	}); err != nil {
		log.Fatal(err)
	}

	// Initialize handler
	h := &handler.Handler{DB: db}

	// Routes
	e.POST("/login", h.Login)
	//=>pokemon
	e.POST("/pokemon", h.CreatePokemon)
	e.GET("/pokemon", h.PokemonsAll)
	e.GET("/pokemon/:id", h.PokemonsFindId)
	e.PUT("/pokemon/:id", h.PokemonsUpdate)
	e.DELETE("/pokemon/:id", h.PokemonsDelete)
	//=>user
	e.POST("/user", h.CreateUser)
	e.GET("/user", h.UserAll)
	//=>upload
	e.GET("/upload", h.Upload)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
