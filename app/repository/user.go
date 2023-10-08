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

type IUserRepository interface {
	CreateUser(ctx context.Context, data dao.User) error
	GetUser(ctx context.Context, filter bson.D) (dao.User, error)
	GetUsers(ctx context.Context, filter bson.D, opt ...*options.FindOptions) ([]dao.User, error)
	UpdateUser(ctx context.Context, filter, updateData bson.D) error
	DeleteUser(ctx context.Context, filter bson.D) error
}

type UserRepository struct {
	db *config.Db
}

func NewUserRepository(db *config.Db) *UserRepository {
	return &UserRepository{db: db}
}

func (r UserRepository) CreateUser(ctx context.Context, data dao.User) error {
	data.Base.ID = primitive.NewObjectID()
	data.Base.CreatedAt = time.Now()
	return Insert(ctx, data, r.db, os.Getenv("DB_NAME"), "Users")
}

func (r UserRepository) GetUser(ctx context.Context, filter bson.D) (dao.User, error) {
	return Get[dao.User](ctx, r.db, os.Getenv("DB_NAME"), "Users", filter)
}

func (r UserRepository) GetUsers(ctx context.Context, filter bson.D, opt ...*options.FindOptions) ([]dao.User, error) {
	return GetList[dao.User](ctx, r.db, os.Getenv("DB_NAME"), "Users", filter, opt...)
}

func (r UserRepository) UpdateUser(ctx context.Context, filter, updateData bson.D) error {
	return Update(ctx, os.Getenv("DB_NAME"), "Users", r.db, filter, updateData)
}

func (r UserRepository) DeleteUser(ctx context.Context, filter bson.D) error {
	return Delete(ctx, os.Getenv("DB_NAME"), "Users", r.db, filter)
}
