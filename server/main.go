package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/middleware"

	socketio "github.com/googollee/go-socket.io"
	"github.com/labstack/echo"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/olebedev/emitter"
)

var db *mongo.Database
var e *echo.Echo
var io *socketio.Server
var emit = emitter.Emitter{}

var instance = "dev-server"

func main() {
	connectDB()
	setupIndexes()

	configureWeb()
	configureSocket()

	startExistingTimers()

	e.Logger.Fatal(e.Start(":8081"))
}

func connectDB() {
	client, err := mongo.NewClient("mongodb://localhost")
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//go testEmit()

	db = client.Database("test")
}

func configureWeb() {
	e = echo.New()
	e.Use(middleware.CORS())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/session/new", newSession)
	e.GET("/session/:uid", getSession)
	e.POST("/session/:uid/time", setTime)
}

func configureSocket() {
	var err error
	io, err = socketio.NewServer([]string{"polling"})
	if err != nil {
		log.Fatal(err)
	}
	e.Any("/socket.io/", echo.WrapHandler(io))
	e.Use(socketioCORS)

	io.SetAllowUpgrades(false)

	io.On("connection", func(so socketio.Socket) {
		log.Println("connection")
		ctx, cancel := context.WithCancel(context.Background())
		so.On("subscribe", func(uid string) {
			log.Println("subscribe")
		S:
			for {
				select {
				case <-ctx.Done():
					break S
				case event := <-emit.Once(uid):
					so.Emit("timeUpdate", int64(event.Float(0)))
				}
			}
		})
		so.On("disconnection", func() {
			cancel()
		})
	})
}

func testEmit() {
	for event := range emit.On("*") {
		fmt.Println(event.OriginalTopic, event.Float(0))
	}
}

func socketioCORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasPrefix(c.Path(), "/socket.io/") {
			if origin := c.Request().Header.Get("Origin"); origin != "" {
				c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
				c.Response().Header().Set("Access-Control-Allow-Origin", origin)
			}
		}
		return next(c)
	}
}
