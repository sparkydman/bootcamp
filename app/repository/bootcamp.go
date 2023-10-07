package repository

import (
	"bootcamp-api/app/model/dao"
	"bootcamp-api/config"
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IBootcampRepository interface {
	AddBootcamp(ctx context.Context, data dao.Bootcamp) error
	GetBootcampById(ctx context.Context, filter bson.M) (dao.Bootcamp, error)
	GetBootcamps(ctx context.Context, filter bson.D, opt ...*options.FindOptions) ([]dao.Bootcamp, error)
	UpdateBootcamp(ctx context.Context, filter, data bson.D) error
	DeleteBootcamp(ctx context.Context, filter bson.D) error
	GetBootcampByFieldName(ctx context.Context, filter bson.D) (dao.Bootcamp, error)
}

type BootcampRepository struct {
	db *config.Db
}

func NewBootcampRepository(db *config.Db) *BootcampRepository {
	return &BootcampRepository{db: db}
}

func (b BootcampRepository) AddBootcamp(ctx context.Context, data dao.Bootcamp) error {
	data.Base.ID = primitive.NewObjectID()
	data.Base.CreatedAt = time.Now()
	_, err := b.db.Client.Database(os.Getenv("DB_NAME")).Collection("Bootcamps").InsertOne(ctx, data)
	return err
}
func (b BootcampRepository) GetBootcampById(ctx context.Context, filter bson.M) (dao.Bootcamp, error) {
	var bootcamp dao.Bootcamp
	if err := b.db.Client.Database(os.Getenv("DB_NAME")).Collection("Bootcamps").FindOne(ctx, filter).Decode(&bootcamp); err != nil {
		return dao.Bootcamp{}, err
	}
	return bootcamp, nil
}
func (b BootcampRepository) GetBootcampByFieldName(ctx context.Context, filter bson.D) (dao.Bootcamp, error) {
	var bootcamp dao.Bootcamp
	if err := b.db.Client.Database(os.Getenv("DB_NAME")).Collection("Bootcamps").FindOne(ctx, filter).Decode(&bootcamp); err != nil {
		return bootcamp, err
	}
	return bootcamp, nil
}
func (b BootcampRepository) GetBootcamps(ctx context.Context, filter bson.D, opt ...*options.FindOptions) ([]dao.Bootcamp, error) {
	var bootcamps []dao.Bootcamp
	cursor, err := b.db.Client.Database(os.Getenv("DB_NAME")).Collection("Bootcamps").Find(ctx, filter, opt...)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &bootcamps); err != nil {
		return nil, err
	}
	return bootcamps, nil
}
func (b BootcampRepository) UpdateBootcamp(ctx context.Context, filter, data bson.D) error {
	_, err := b.db.Client.Database(os.Getenv("DB_NAME")).Collection("Bootcamps").UpdateOne(ctx, filter, data)
	if err != nil {
		return err
	}
	return nil
}
func (b BootcampRepository) DeleteBootcamp(ctx context.Context, filter bson.D) error {
	_, err := b.db.Client.Database(os.Getenv("DB_NAME")).Collection("Bootcamps").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
