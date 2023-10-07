package config

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Db struct {
	*mongo.Client
}

func NewConnection(uri string) (*Db, error) {
	c, err := connect(uri)
	if err != nil {
		return nil, err
	}
	return &Db{c}, nil
}

func (d *Db) IsConnected() bool {
	var result bson.M
	if err := d.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return false
	}
	return true
}

func (d *Db) Disconnect() func() {
	return func() {
		if err := d.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}

func connect(uri string) (*mongo.Client, error) {
	version := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(version)

	return mongo.Connect(context.TODO(), opts)
}
