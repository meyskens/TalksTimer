package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/labstack/echo/middleware"

	socketio "github.com/googollee/go-socket.io"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/olebedev/emitter"
)

type config struct {
	Port     int    `default:"8081"`
	Instance string `default:"dev-server"`
}

var conf = config{}

var db *mongo.Database
var e *echo.Echo
var io *socketio.Server
var emit = emitter.Emitter{}

func main() {
	err := envconfig.Process("talkstimer", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	connectDB()
	setupIndexes()

	configureWeb()
	configureSocket()

	startExistingTimers()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Port)))
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
	e.Static("/", "www")
	e.POST("/session/new", newSession)
	e.GET("/session/:uid", getSession)
	e.POST("/session/:uid/time", setTime)
	e.POST("/session/:uid/message", sendMessage)
	e.GET("/session/:uid/colors", getColors)
	e.POST("/session/:uid/colors", setColors)
}

func configureSocket() {
	var err error
	io, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	e.Any("/socket.io/", echo.WrapHandler(io))
	e.Use(socketioCORS)

	io.SetPingInterval(time.Second)
	io.On("connection", func(so socketio.Socket) {
		log.Println("connection")
		so.On("disconnection", func() {
			log.Println("disconnection")
		})

		so.On("subscribe", func(uid string) {
			log.Println("subscribe", uid)
			so.Join(uid)
		})
	})
	io.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
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
