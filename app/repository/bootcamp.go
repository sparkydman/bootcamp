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
	GetBootcamp(ctx context.Context, filter bson.D) (dao.Bootcamp, error)
	GetBootcamps(ctx context.Context, filter bson.D, opt ...*options.FindOptions) ([]dao.Bootcamp, error)
	UpdateBootcamp(ctx context.Context, filter, data bson.D) error
	DeleteBootcamp(ctx context.Context, filter bson.D) error
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
	return Insert(ctx, data, b.db, os.Getenv("DB_NAME"), "Bootcamps")
}

func (b BootcampRepository) GetBootcamp(ctx context.Context, filter bson.D) (dao.Bootcamp, error) {
	return Get[dao.Bootcamp](ctx, b.db, os.Getenv("DB_NAME"), "Bootcamps", filter)
}

func (b BootcampRepository) GetBootcamps(ctx context.Context, filter bson.D, opts ...*options.FindOptions) ([]dao.Bootcamp, error) {
	return GetList[dao.Bootcamp](ctx, b.db, os.Getenv("DB_NAME"), "Bootcamps", filter, opts...)
}

func (b BootcampRepository) UpdateBootcamp(ctx context.Context, filter, data bson.D) error {
	return Update(ctx, os.Getenv("DB_NAME"), "Bootcamps", b.db, filter, data)
}

func (b BootcampRepository) DeleteBootcamp(ctx context.Context, filter bson.D) error {
	return Delete(ctx, os.Getenv("DB_NAME"), "Bootcamps", b.db, filter)
}
