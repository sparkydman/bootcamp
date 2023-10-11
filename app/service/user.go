package service

import (
	"bootcamp-api/app/model/dao"
	"bootcamp-api/app/model/dto"
	"bootcamp-api/app/repository"
	"bootcamp-api/config"
	"bootcamp-api/utils"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	logger "github.com/sirupsen/logrus"
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
	GetLoggedInUser(c *gin.Context)
	GetToken(c *gin.Context)
	Logout(c *gin.Context)
}

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s UserService) CreateUser(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)

	var request dto.CreateUserRequest
	code := utils.BadRequestErrorCode
	if err := c.ShouldBindJSON(&request); err != nil {
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}

	if err := request.Vaildate(); err != nil {
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}

	_, err := s.repo.GetUser(context.TODO(), bson.D{{"email", request.Email}})
	if err == nil {
		logger.Error("User already exists")
		code = utils.ConflictErrorCode
		code.SetMessage("user already exists")
		utils.PanicException(code)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.PanicException(utils.ServerErrorCode)
	}
	request.Password = string(hashedPassword)

	user, err := utils.EmbedStructFlat[dao.User](request)
	if err != nil {
		logger.Error("couldn't embed struct", err)
		utils.PanicException(utils.ServerErrorCode)
	}
	if err := s.repo.CreateUser(context.TODO(), user); err != nil {
		utils.PanicException(utils.ServerErrorCode)
	}

	c.JSON(http.StatusCreated, utils.SetResponse(true, utils.SuccessfulCode, utils.NULL()))
}

func (s UserService) LoginUser(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)

	var request dto.LoginUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.PanicException(utils.CredentialsErrorCode)
	}

	user, err := s.repo.GetUser(context.TODO(), bson.D{{"email", request.Email}})
	if err != nil {
		logger.Error("failed to get user by email: %s, error: %v", request.Email, err)
		utils.PanicException(utils.CredentialsErrorCode)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		utils.PanicException(utils.CredentialsErrorCode)
	}

	payload := user
	payload.Password = ""
	accessToken, err := generateAccessToken(payload, time.Now().Add(30*time.Minute))
	if err != nil {
		logger.Error("failed to generate access token", err)
		code := utils.UnAuthorizedErrorCode
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}

	refreshToken, err := generateRefreshToken(user.ID, time.Now().Add(30*24*time.Hour))
	if err != nil {
		logger.Error("failed to generate refresh token", err)
		code := utils.UnAuthorizedErrorCode
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}

	loginUser := dto.LoginUserResponse{User: user, AccessToken: accessToken, RefreshToken: refreshToken}

	loginUser.Password = ""
	c.SetCookie("refresh_token", refreshToken, 60*60*24*30, "/", "", false, false)
	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, loginUser))
}

func (s UserService) GetUserById(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)

	param := c.Params.ByName("id")
	id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		logger.Errorf("failed to create object ID: %s, Error: %v", param, err)
		code := utils.BadRequestErrorCode
		code.SetMessage("Invalid user id")
		utils.PanicException(code)
	}

	filter := bson.D{{"_id", id}}
	user, err := s.repo.GetUser(context.TODO(), filter)
	if err != nil {
		logger.Errorf("failed to get user by id: %v, error: %v", id, err)
		utils.PanicException(utils.ServerErrorCode)
	}

	user.Password = ""
	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, user))
}

func (s UserService) GetUsers(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)

	filter := bson.D{}
	opts := utils.FindOptions(c)
	opts = append(opts, options.Find().SetProjection(bson.D{{"password", 0}}))

	users, err := s.repo.GetUsers(context.TODO(), filter, opts...)
	if err != nil {
		logger.Errorf("error while getting users %v", err)
		utils.PanicException(utils.ServerErrorCode)
	}

	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, users))
}

func (s UserService) UpdateUser(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)
	var err error

	var user dto.UserRequest
	code := utils.BadRequestErrorCode
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("failed to bind JSON", err)
		utils.PanicException(code)
	}

	u := c.MustGet("LoggedInUser").(dao.User)
	var id primitive.ObjectID
	if u.Role == "admin" {
		id, err = primitive.ObjectIDFromHex(c.Params.ByName("userid"))
		if err != nil {
			logger.Error("failed to validate request parameter", err)
			code.SetMessage("Invalid user id")
			utils.PanicException(code)
		}
	}
	id = u.ID
	if _, err := s.repo.GetUser(context.TODO(), bson.D{{"_id", id}}); err != nil {
		logger.Errorf("failed to get user by id: %v, error: %v", id, err)
		errorCode := utils.NotFoundErrorCode
		errorCode.SetMessage("user not found")
		utils.PanicException(errorCode)
	}

	if err := user.Vaildate(u.Role); err != nil {
		logger.Error("failed to validate request body", err)
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}

	filter := bson.D{{"_id", id}}
	data := bson.D{{"$set", bson.D{{"name", user.Name}, {"role", user.Role}, {"updated_at", time.Now()}}}}

	if err := s.repo.UpdateUser(context.TODO(), filter, data); err != nil {
		logger.Error("failed to update user", err)
		utils.PanicException(utils.ServerErrorCode)
	}

	if u.Role != "admin" {
		u.Role = user.Role
		u.UpdatedAt = time.Now()
		u.Name = user.Name

		accessToken, err := generateAccessToken(u, time.Now().Add(30*time.Minute))
		if err != nil {
			logger.Error("failed to generate access token", err)
			code := utils.ServerErrorCode
			utils.PanicException(code)
		}

		c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, struct {
			dao.User
			AccessToken string `json:"access_token"`
		}{
			u,
			accessToken,
		}))
	} else {
		c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, utils.NULL()))
	}
}

func (s UserService) DeleteUser(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)
	var err error

	user := c.MustGet("LoggedInUser").(dao.User)
	param := c.Params.ByName("userid")
	var id primitive.ObjectID
	if user.Role == "admin" {
		id, err = primitive.ObjectIDFromHex(param)
		if err != nil {
			logger.Errorf("failed to validate user id: %s, error: %v", param, err)
			code := utils.BadRequestErrorCode
			code.SetMessage("Invalid user ID")
			utils.PanicException(code)
		}
	}
	id = user.ID

	if err := s.repo.DeleteUser(context.TODO(), bson.D{{"_id", id}}); err != nil {
		logger.Errorf("failed to delete user with id: %s, error: %v", param, err)
		utils.PanicException(utils.ServerErrorCode)
	}
	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, utils.NULL()))
}

func (s UserService) GetLoggedInUser(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)
	user := c.MustGet("LoggedInUser")

	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, user))
}

func (s UserService) GetToken(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)
	code := utils.BadRequestErrorCode

	cookies, err := c.Cookie("refresh_token")
	if err != nil {
		logger.Errorf("failed to retrieve refresh token from cookies, error: %v", err)
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}
	decode, err := config.VerifyToken[primitive.ObjectID](cookies, []byte(os.Getenv("REFRESH_TOKEN_KEY")), primitive.ObjectID{})
	if err != nil {
		logger.Errorf("failed to decode refresh token, error: %v", err)
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}

	user, err := s.repo.GetUser(context.TODO(), bson.D{{"_id", decode}})
	if err != nil {
		logger.Errorf("failed to retrieve user from db, error: %v", err)
		code = utils.NotFoundErrorCode
		code.SetMessage("user not found")
		utils.PanicException(code)
	}

	user.Password = ""
	newToken, err := generateAccessToken(user, time.Now().Add(30*time.Minute))
	if err != nil {
		logger.Errorf("failed to generate new access token, error: %v", err)
		code = utils.ServerErrorCode
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}

	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, struct {
		AccessToken string `json:"access_token,omitempty"`
	}{AccessToken: newToken}))
}

func (s UserService) Logout(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)
	errorCode := utils.ServerErrorCode

	user := c.MustGet("LoggedInUser").(dao.User)
	refreshToken, err := generateRefreshToken(user.ID, time.Now())
	if err != nil {
		utils.PanicException(errorCode)
	}

	c.SetCookie("refresh_token", refreshToken, 0, "/", "", false, false)
	c.Set("LoggedInUser", nil)
	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, struct {
		Message string `json:"message"`
	}{Message: "Logout successfully"}))
}

func generateAccessToken(payload dao.User, age time.Time) (string, error) {
	accessTokenClaims := config.TokenClaims[dao.User]{
		payload,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(age),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    os.Getenv("TOKEN_ISSUER"),
			Subject:   os.Getenv("TOKEN_SUBJECT"),
		},
	}

	return accessTokenClaims.GenerateToken([]byte(os.Getenv("ACCESS_TOKEN_KEY")))
}

func generateRefreshToken(id primitive.ObjectID, age time.Time) (string, error) {
	refreshTokenClaims := config.TokenClaims[primitive.ObjectID]{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(age),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    os.Getenv("TOKEN_ISSUER"),
			Subject:   os.Getenv("TOKEN_SUBJECT"),
		},
	}

	return refreshTokenClaims.GenerateToken([]byte(os.Getenv("REFRESH_TOKEN_KEY")))
}
