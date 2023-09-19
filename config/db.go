package config

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Db struct {
	*mongo.Client
	Uri string
}

func (d *Db) IsConnected() bool {
	var result bson.M
	if err := d.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return false
	}
	return true
}

func (d *Db) Connect() error {
	version := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(d.Uri).SetServerAPIOptions(version)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	d.Client = client
	return nil
}

func (d *Db) Disconnect() func() {
	return func() {
		if err := d.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}
}
