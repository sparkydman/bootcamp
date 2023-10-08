package repository

import (
	"bootcamp-api/config"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Get[T any](ctx context.Context, db *config.Db, dbname string, collection string, filter bson.D, opts ...*options.FindOneOptions)(T, error){
	var data T
	if err := db.Client.Database(dbname).Collection(collection).FindOne(ctx, filter, opts...).Decode(&data); err != nil {
		return data, err
	}
	return data, nil
}

func GetList[T any](ctx context.Context, db *config.Db, dbname string, collection string, filter bson.D, opts ...*options.FindOptions) ([]T, error) {
	var data []T
	cursor, err := db.Client.Database(dbname).Collection(collection).Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func Insert[T any](ctx context.Context, data T, db *config.Db, dbname string, collection string, opts ...*options.InsertOneOptions) error {
	_, err := db.Client.Database(dbname).Collection(collection).InsertOne(ctx, data, opts...)

	return err
}

func InsertMany[T any](ctx context.Context, data []T, db *config.Db, dbname string, collection string, opts ...*options.InsertManyOptions) error {
	var payload []interface{}
	for _, d := range data{
		payload = append(payload, d)
	}
	_, err := db.Client.Database(dbname).Collection(collection).InsertMany(ctx, payload, opts...)

	return err
}

func Update(ctx context.Context, dbname string, collection string, db *config.Db, filter, updateData bson.D, opts ...*options.UpdateOptions) error {
	_, err := db.Client.Database(dbname).Collection(collection).UpdateOne(ctx, filter, updateData, opts...)
	if err != nil {
		return err
	}
	return nil
}

func UpdateMany(ctx context.Context, dbname string, collection string, db *config.Db, filter, updateData bson.D, opts ...*options.UpdateOptions) error {
	_, err := db.Client.Database(dbname).Collection(collection).UpdateMany(ctx, filter, updateData, opts...)
	if err != nil {
		return err
	}
	return nil
}

func Delete(ctx context.Context, dbname string, collection string, db *config.Db, filter bson.D, opts ...*options.DeleteOptions) error {
	_, err := db.Client.Database(dbname).Collection(collection).DeleteOne(ctx, filter, opts...)
	if err != nil {
		return err
	}
	return nil
}

func DeleteMany(ctx context.Context, dbname string, collection string, db *config.Db, filter bson.D, opts ...*options.DeleteOptions) error {
	_, err := db.Client.Database(dbname).Collection(collection).DeleteMany(ctx, filter, opts...)
	if err != nil {
		return err
	}
	return nil
}