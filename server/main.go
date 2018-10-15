package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

	e.Logger.Fatal(e.Start(":8080"))
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

	go testEmit()

	db = client.Database("test")
}

func configureWeb() {
	e = echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/session/new", newSession)
	e.GET("/session/:uid", getSession)
	e.POST("/session/:uid/time", setTime)
}

func configureSocket() {
	var err error
	io, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	e.Any("/socket.io/", echo.WrapHandler(io))

	io.On("connection", func(so socketio.Socket) {
		ctx, cancel := context.WithCancel(context.Background())
		so.On("subscribe", func(uid string) {
		S:
			for {
				select {
				case <-ctx.Done():
					break S
				case event := <-emit.On(uid):
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
