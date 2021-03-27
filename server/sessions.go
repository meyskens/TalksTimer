package main

import (
	"context"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mongodb/mongo-go-driver/bson"
)

// Session keeps the info of a timer session
type Session struct {
	ID          objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Key         string            `bson:"key" json:"key"`
	SecondsLeft int64             `bson:"secondsLeft" json:"secondsLeft"` // seconds
	Created     time.Time         `bson:"created" json:"created"`
	Instance    string            `bson:"instance" json:"instance"` // this links a timer to a server instance
}

func newSession(c echo.Context) error {
	uid, err := uuid.NewRandom()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create unique key"})
	}
	session := Session{
		Key:      uid.String(),
		Instance: conf.Instance,
		Created:  time.Now(),
	}
	_, err = db.Collection("sessions").InsertOne(context.Background(), session)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create database record"})
	}
	return c.JSON(http.StatusOK, session)
}

func getSession(c echo.Context) error {
	uid := c.Param("uid")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no uid given"})
	}

	session := Session{}
	err := db.Collection("sessions").FindOne(context.Background(), bson.NewDocument(bson.EC.String("key", uid))).Decode(&session)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "did not find session"})
	}

	return c.JSON(http.StatusOK, session)
}

type setTimeBody struct {
	Seconds interface{} `json:"seconds" form:"seconds" query:"seconds"`
}

func setTime(c echo.Context) error {
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

	body := setTimeBody{}
	c.Bind(&body)

	var seconds int64

	if reflect.TypeOf(body.Seconds).Kind() == reflect.Int {
		seconds = body.Seconds.(int64)
	}
	if reflect.TypeOf(body.Seconds).Kind() == reflect.Float64 {
		seconds = int64(body.Seconds.(float64))
	}
	if reflect.TypeOf(body.Seconds).Kind() == reflect.String {
		seconds, _ = strconv.ParseInt(body.Seconds.(string), 10, 64)
	}

	_, err = db.Collection("sessions").UpdateOne(context.Background(), bson.NewDocument(bson.EC.String("key", uid)), bson.NewDocument(bson.EC.SubDocumentFromElements("$set", bson.EC.Int64("secondsLeft", seconds))))
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error setting time left"})
	}

	if cancel, exists := timerCancels[uid]; exists {
		cancel()
	}
	startTimer(uid, seconds)

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
