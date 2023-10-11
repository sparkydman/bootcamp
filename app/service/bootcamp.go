package service

import (
	"bootcamp-api/app/model/dao"
	"bootcamp-api/app/model/dto"
	"bootcamp-api/app/repository"
	"bootcamp-api/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IBootcampService interface {
	AddBootcamp(c *gin.Context)
	GetBootcampsByCreator(c *gin.Context)
}

type Bootcamp struct {
	repo repository.IBootcampRepository
}

func NewBootcampService(repo repository.IBootcampRepository) *Bootcamp {
	return &Bootcamp{repo: repo}
}

func (b Bootcamp) AddBootcamp(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)

	var request dto.BootcampRequest
	code := utils.BadRequestErrorCode
	//bind request body
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("failed to bind request", err)
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}
	//vailable request body
	if err := request.Validate(dto.ValidCareers()); err != nil {
		logger.Error("failed to validate request", err)
		code.SetMessage(err.Error())
		utils.PanicException(code)
	}
	request.Slugify()

	creatorId, _ := primitive.ObjectIDFromHex(request.CreatedBy)
	//check if bootcamp is already exist for this user
	_, err := b.repo.GetBootcamp(context.TODO(), primitive.D{{"slug", request.Slug}, {"created_by", creatorId}})
	if err == nil {
		errorMsg := fmt.Sprintf("bootcamp already exists for user - %s", request.CreatedBy)
		logger.Error(errorMsg)
		code = utils.ConflictErrorCode
		code.SetMessage(errorMsg)
		utils.PanicException(code)
	}

	//add bootcamp
	payload, err := utils.EmbedStructFlat[dao.Bootcamp](request)
	if err != nil {
		logger.Error("failed to embed request", err)
		code = utils.ServerErrorCode
		utils.PanicException(code)
	}
	if err := b.repo.AddBootcamp(context.TODO(), payload); err != nil {
		logger.Error("failed to create bootcamp", err)
		code = utils.ServerErrorCode
		utils.PanicException(code)
	}

	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, utils.NULL()))
}

func (b Bootcamp) GetBootcampsByCreator(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)
	errorCode := utils.BadRequestErrorCode

	creator := c.Params.ByName("creator")
	creatorId, err := primitive.ObjectIDFromHex(creator)
	if err != nil {
		logger.Error("failed to validate mongodb id", err)
		errorCode.SetMessage(err.Error())
		utils.PanicException(errorCode)
	}

	filter := bson.D{{"created_by", creatorId}}
	bootcamps, err := b.repo.GetBootcamps(context.TODO(), filter)
	if err != nil {
		logger.Error("failed to get bootcamps from a creator", err)
		errorCode = utils.ServerErrorCode
		utils.PanicException(errorCode)
	}

	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, bootcamps))
}

func (b Bootcamp) GetBootcamps(c *gin.Context) {
	defer utils.ResponseErrorHandler(c)
	errorCode := utils.BadRequestErrorCode

	creator := c.Params.ByName("creator")
	creatorId, err := primitive.ObjectIDFromHex(creator)
	if err != nil {
		logger.Error("failed to validate mongodb id", err)
		errorCode.SetMessage(err.Error())
		utils.PanicException(errorCode)
	}

	filter := bson.D{{"created_by", creatorId}}
	bootcamps, err := b.repo.GetBootcamps(context.TODO(), filter)
	if err != nil {
		logger.Error("failed to get bootcamps from a creator", err)
		errorCode = utils.ServerErrorCode
		utils.PanicException(errorCode)
	}

	c.JSON(http.StatusOK, utils.SetResponse(true, utils.SuccessfulCode, bootcamps))
}
