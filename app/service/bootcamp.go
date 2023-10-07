package service

import (
	"bootcamp-api/app/model/dao"
	"bootcamp-api/app/model/dto"
	"bootcamp-api/app/repository"
	"bootcamp-api/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	logger "github.com/sirupsen/logrus"
)

type IBootcampService interface {
	AddBootcamp(c *gin.Context)
}

type bootcamp struct {
	repo repository.IBootcampRepository
}

func NewBootcampService(repo repository.IBootcampRepository) IBootcampService {
	return &bootcamp{repo: repo}
}

func (b bootcamp) AddBootcamp(c *gin.Context){
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

	//check if bootcamp is already exist
	_, err := b.repo.GetBootcampByFieldName(context.TODO(), primitive.D{{"slug", request.Slug}})
	if err == nil {
		logger.Error("bootcamp already exists")
		code = utils.ConflictErrorCode
		code.SetMessage("bootcamp already exists")
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