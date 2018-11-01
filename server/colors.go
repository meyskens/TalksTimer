package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/labstack/echo"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Colors keeps the color options for for a session
type Colors struct {
	ID         objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	SessionKey string            `bson:"sessionKey" json:"sessionKey"`
	Created    time.Time         `bson:"created" json:"created"`
	Options    []ColorOptions    `bson:"options" json:"options"`
}

// ColorOptions keep the color for a certain time
type ColorOptions struct {
	Color string `bson:"color" json:"color"`
	From  int64  `bson:"from" json:"from"`
}

func setColors(c echo.Context) error {
	uid := c.Param("uid")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no uid given"})
	}

	colors := Colors{}
	c.Bind(&colors)
	spew.Dump(colors)
	colors.ID = objectid.NilObjectID
	colors.SessionKey = uid
	colors.Created = time.Now()

	db.Collection("colors").DeleteMany(context.Background(), bson.NewDocument(bson.EC.String("sessionKey", uid))) // delete old ones
	_, err := db.Collection("colors").InsertOne(context.Background(), colors)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not create database record"})
	}

	io.BroadcastTo(uid, "newColors", true)

	return c.JSON(http.StatusOK, colors)
}

func getColors(c echo.Context) error {
	uid := c.Param("uid")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "no uid given"})
	}

	colors := Colors{}
	db.Collection("colors").FindOne(context.Background(), bson.NewDocument(bson.EC.String("sessionKey", uid))).Decode(&colors)

	return c.JSON(http.StatusOK, colors)
}
