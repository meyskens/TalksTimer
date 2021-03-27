package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Messages keeps the message options for for a session
type Messages struct {
	ID         objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	SessionKey string            `bson:"sessionKey" json:"sessionKey"`
	Created    time.Time         `bson:"created" json:"created"`
	Messages   []Message         `bson:"messages" json:"messages"`
}

// Message keeps a string of a message with a given color
type Message struct {
	Message string `bson:"message" json:"message"`
	Color   string `bson:"color" json:"color"`
}

var defaultMessages = []Message{
	Message{
		Message: "Please repeat the question",
		Color:   "#ffffff",
	},
	Message{
		Message: "Hold the mic closer",
		Color:   "#ffffff",
	},
}

func setMessages(c echo.Context) error {
	uid := c.Param("uid")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no uid given"})
	}

	messages := Messages{}
	c.Bind(&messages)

	messages.ID = objectid.NilObjectID
	messages.SessionKey = uid
	messages.Created = time.Now()

	db.Collection("messages").DeleteMany(context.Background(), bson.NewDocument(bson.EC.String("sessionKey", uid))) // delete old ones
	_, err := db.Collection("messages").InsertOne(context.Background(), messages)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create database record"})
	}

	return c.JSON(http.StatusOK, messages)
}

func getMessages(c echo.Context) error {
	uid := c.Param("uid")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no uid given"})
	}

	messages := Messages{}
	db.Collection("messages").FindOne(context.Background(), bson.NewDocument(bson.EC.String("sessionKey", uid))).Decode(&messages)

	if len(messages.Messages) == 0 {
		messages.Messages = defaultMessages
	}

	return c.JSON(http.StatusOK, messages)
}

type sendMessageBody struct {
	Message string `json:"message" form:"message" query:"message"`
	Color   string `json:"color" form:"color" query:"color"`
}

func sendMessage(c echo.Context) error {
	uid := c.Param("uid")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no uid given"})
	}

	count, err := db.Collection("sessions").Count(context.Background(), bson.NewDocument(bson.EC.String("key", uid)))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error looking up session"})
	}
	if count < 1 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "session not found"})
	}

	body := sendMessageBody{}
	c.Bind(&body)

	io.BroadcastTo(uid, "message", body)

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
