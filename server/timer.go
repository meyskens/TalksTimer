package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
)

var timerContexts = map[string]context.Context{}
var timerCancels = map[string]context.CancelFunc{}

func startExistingTimers() {
	c, err := db.Collection("sessions").Find(context.Background(), bson.NewDocument(bson.EC.String("instance", conf.Instance)))
	if err != nil {
		log.Println(err)
		return
	}

	for c.Next(context.Background()) {
		session := Session{}
		err := c.Decode(&session)
		if err != nil {
			log.Println(err)
			continue
		}
		if session.SecondsLeft > 0 {
			startTimer(session.Key, session.SecondsLeft)
		}
	}
}

func startTimer(key string, seconds int64) {
	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", seconds))
	timerContexts[key], timerCancels[key] = context.WithTimeout(context.Background(), duration)
	go runTimer(key, duration)
}

func runTimer(key string, duration time.Duration) {
	ctx := timerContexts[key]
	ticker := time.NewTicker(time.Second)
L:
	for {
		select {
		case <-ticker.C:
			duration = duration - time.Second
			go updateDuration(key, duration)
			if duration.Seconds() == 0 {
				break L
			}
		case <-ctx.Done():
			duration = 0 * time.Second
			go updateDuration(key, duration)
			break L
		}
	}
}

func updateDuration(key string, duration time.Duration) {
	emit.Emit(key, duration.Seconds())
	io.BroadcastTo(key, "timeUpdate", duration.Seconds())
	_, err := db.Collection("sessions").UpdateOne(context.Background(), bson.NewDocument(bson.EC.String("key", key)), bson.NewDocument(bson.EC.SubDocumentFromElements("$set", bson.EC.Int64("secondsLeft", int64(duration.Seconds())))))
	if err != nil {
		log.Println(err)
	}
}
