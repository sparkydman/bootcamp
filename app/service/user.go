package service

import (
	"bootcamp-api/app/model/dao"
	"bootcamp-api/app/model/dto"
	"bootcamp-api/app/repository"
	"bootcamp-api/config"
	"bootcamp-api/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	CreateUser(c *gin.Context)
	GetUserById(c *gin.Context)
	LoginUser(c *gin.Context)
	GetUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) IUserService {
	return &userService{repo: repo}
}

func (s userService) CreateUser(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)

	var request dto.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.PanicException(utils.BadRequestErrorCode)
	}
	
	code := utils.BadRequestErrorCode
	if err := request.Vaildate(); err != nil {
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}

	_, err := s.repo.GetUserByFieldName(context.TODO(), bson.D{{"email", request.Email}})
	if err == nil{
		code.SetMessage("user already exists")
		utils.PanicException(code)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.PanicException(utils.ServerErrorCode)
	}
	request.Password = string(hashedPassword)

	user := dao.User{
		Name: request.Name,
		Email: request.Email,
		Role: request.Role,
		Password: request.Password,
	}
	if err := s.repo.CreateUser(context.TODO(), user); err != nil {
		utils.PanicException(utils.ServerErrorCode)
	}

	c.JSON(http.StatusCreated, utils.SetResponse(true, utils.SuccessfulCode, utils.NULL()))
}

func (s userService) LoginUser(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)

	var request dto.LoginUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.PanicException(utils.CredentialsErrorCode)
	}

	user, err := s.repo.GetUserByFieldName(context.TODO(), bson.D{{"email", request.Email}})
	if err != nil {
		utils.PanicException(utils.CredentialsErrorCode)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		utils.PanicException(utils.CredentialsErrorCode)
	}

	accessTokenClaims := config.TokenClaims[dao.User]{
		dao.User{
			Email: user.Email,
			Name: user.Name,
			Role: user.Role,
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer: os.Getenv("TOKEN_ISSUER"),
			Subject: os.Getenv("TOKEN_SUBJECT"),
		},
	}

	accessToken, err := accessTokenClaims.GenerateToken([]byte(os.Getenv("ACCESS_TOKEN_KEY")))
	if err != nil {
		code := utils.UnAuthorizedErrorCode
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}

	age := time.Now().Add(30 * 24 * time.Hour);
	refreshTokenClaims := config.TokenClaims[dao.User]{
		dao.User{
			Email: user.Email,
			Name: user.Name,
			Role: user.Role,
		},
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(age),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer: os.Getenv("TOKEN_ISSUER"),
			Subject: os.Getenv("TOKEN_SUBJECT"),
		},
	}

	refreshToken, err := refreshTokenClaims.GenerateToken([]byte(os.Getenv("REFRESH_TOKEN_KEY")))
	if err != nil {
		code := utils.UnAuthorizedErrorCode
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}

	loginUser := dto.LoginUserResponse{User: user, AccessToken: accessToken, RefreshToken: refreshToken}

	loginUser.Password = ""
	c.SetCookie("refresh_token", refreshToken, 60 * 60 * 24 * 30, "/", "", false, false)
	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, loginUser))
}

func (s userService) GetUserById(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)

	id, err := primitive.ObjectIDFromHex(c.Params.ByName("id"))
	if err != nil {
		code := utils.BadRequestErrorCode
		code.SetMessage("Invalid user id")
		utils.PanicException(code)
	}

	filter := bson.M{"_id": id}
	user, err := s.repo.GetUserById(context.TODO(), filter)
	if err != nil {
		utils.PanicException(utils.ServerErrorCode)
	}

	user.Password = ""
	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, user))
}

func (s userService) GetUsers(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)
	// fmt.Println("url from middleware", c.MustGet("RequestUrl").(string))
	filter := bson.D{}
	querySort := c.Query("sort")
	querySelect := c.Query("select")
	opts := []*options.FindOptions{
		options.Find().SetProjection(bson.D{{"password", 0}}),
	}
	if querySort != "" {
		opts = append(opts, options.Find().SetSort(bson.D{{querySort, 1}}))
	}
	if querySelect != "" {
		fields := strings.Split(querySelect, ",")
		var s bson.D
		for _, field := range fields {
			s = append(s, bson.E{Key: field, Value: 1})
		}
		opts = append(opts, options.Find().SetProjection(s))
	}

	users, err := s.repo.GetUsers(context.TODO(), filter, opts...)
	if err != nil {
		fmt.Errorf("error while getting users %v\n", err)
		utils.PanicException(utils.ServerErrorCode)
	}

	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, users))
}

func (s userService) UpdateUser(c *gin.Context){
	defer utils.ResponseErrorHandler(c)

	var user dto.UserRequest
	code := utils.BadRequestErrorCode
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("failed to bind JSON", err)
		utils.PanicException(code)
	}

	id := c.Params.ByName("userid")
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("failed to validate request parameter", err)
		code.SetMessage("Invalid user id")
		utils.PanicException(code)
	}
	if _, err := s.repo.GetUserById(context.TODO(), bson.M{"_id": docId}); err != nil {
		errorCode := utils.NotFoundErrorCode
		errorCode.SetMessage("user not found")
		utils.PanicException(errorCode)
	}

	if err := user.Vaildate(); err != nil {
		log.Println("failed to validate request body", err)
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}

	filter := bson.D{{"_id", docId}}
	data := bson.D{{"$set", bson.D{{"name", user.Name},{"role", user.Role}, {"updated_at", time.Now()}}}}
	
	if err := s.repo.UpdateUser(context.TODO(), filter, data); err != nil{
		log.Println("failed to update user", err)
		utils.PanicException(utils.ServerErrorCode)
	}
	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, utils.NULL()))
}

func (s userService) DeleteUser(c *gin.Context){
	defer utils.ResponseErrorHandler(c)

	id := c.Params.ByName("userid")
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		code := utils.BadRequestErrorCode
		code.SetMessage("Invalid user ID")
		utils.PanicException(code)
	}

	if err := s.repo.DeleteUser(context.TODO(), bson.D{{"_id", docId}}); err != nil{
		utils.PanicException(utils.ServerErrorCode)
	}
	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, utils.NULL()))
}