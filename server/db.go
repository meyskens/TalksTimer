package main

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
)

func setupIndexes() {
	options := mongo.NewIndexOptionsBuilder()
	options.Background(true)
	EnsureIndex(context.Background(), db.Collection("sessions"), []string{"key"}, options)
	EnsureIndex(context.Background(), db.Collection("sessions"), []string{"key", "instance"}, options)
	EnsureIndex(context.Background(), db.Collection("colors"), []string{"key"}, options)

	expireOpts := mongo.NewIndexOptionsBuilder()
	expireOpts.Background(true)
	expireOpts.ExpireAfterSeconds(2592000) //30 days
	EnsureIndex(context.Background(), db.Collection("sessions"), []string{"created"}, expireOpts)
	EnsureIndex(context.Background(), db.Collection("colors"), []string{"created"}, expireOpts)
}

// Thanks to https://gist.github.com/bweston92/5a796e15a6d7f436755795018dea9c1a

// EnsureIndex will ensure the index model provided is on the given collection.
func EnsureIndex(ctx context.Context, c *mongo.Collection, keys []string, opts *mongo.IndexOptionsBuilder) error {
	idxs := c.Indexes()

	ks := bson.NewDocument()
	for _, k := range keys {
		// todo - add support for sorting index.
		ks.Append(bson.EC.Int64(k, -1))
	}
	idm := mongo.IndexModel{
		Keys:    ks,
		Options: opts.Build(),
	}

	v := idm.Options.Lookup("name")
	if v == nil {
		return errors.New("must provide a key name for index")
	}
	expectedName := v.StringValue()

	cur, err := idxs.List(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to list indexes")
	}

	found := false
	for cur.Next(ctx) {
		d := bson.NewDocument()

		if err := cur.Decode(d); err != nil {
			return errors.Wrap(err, "unable to decode bson index document")
		}

		v := d.Lookup("name")
		if v != nil && v.StringValue() == expectedName {
			found = true
			break
		}
	}

	if found {
		return nil
	}

	_, err = idxs.CreateOne(ctx, idm)
	return err
}
