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
	GetUserById(ctx context.Context, filter bson.M) (dao.User, error)
	GetUserByFieldName(ctx context.Context, filter bson.D) (dao.User, error)
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
	_, err := r.db.Client.Database(os.Getenv("DB_NAME")).Collection("Users").InsertOne(ctx, data)

	return err
}

func (r UserRepository) GetUserById(ctx context.Context, filter bson.M) (dao.User, error) {
	var user dao.User
	if err := r.db.Client.Database(os.Getenv("DB_NAME")).Collection("Users").FindOne(ctx, filter).Decode(&user); err != nil {
		return user, err
	}
	return user, nil
}
func (r UserRepository) GetUserByFieldName(ctx context.Context, filter bson.D) (dao.User, error) {
	var user dao.User
	if err := r.db.Client.Database(os.Getenv("DB_NAME")).Collection("Users").FindOne(ctx, filter).Decode(&user); err != nil {
		return user, err
	}
	return user, nil
}
func (r UserRepository) GetUsers(ctx context.Context, filter bson.D, opt ...*options.FindOptions) ([]dao.User, error) {
	var users []dao.User
	cursor, err := r.db.Client.Database(os.Getenv("DB_NAME")).Collection("Users").Find(ctx, filter, opt...)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
func (r UserRepository) UpdateUser(ctx context.Context, filter, updateData bson.D) error {
	_, err := r.db.Client.Database(os.Getenv("DB_NAME")).Collection("Users").UpdateOne(ctx, filter, updateData)
	if err != nil {
		return err
	}
	return nil
}
func (r UserRepository) DeleteUser(ctx context.Context, filter bson.D) error {
	_, err := r.db.Client.Database(os.Getenv("DB_NAME")).Collection("Users").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
