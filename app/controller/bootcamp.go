package controller

import (
	"bootcamp-api/app/service"

	"github.com/gin-gonic/gin"
)

type IBootcampController interface {
	AddBootcamp(c *gin.Context)
	GetBootcampsByCreator(c *gin.Context)
}

type BootcampController struct {
	svc service.IBootcampService
}

func NewBootcampController(svc service.IBootcampService) *BootcampController {
	return &BootcampController{svc: svc}
}
func (b BootcampController) AddBootcamp(c *gin.Context) {
	b.svc.AddBootcamp(c)
}
func (b BootcampController) GetBootcampsByCreator(c *gin.Context){
	b.svc.GetBootcampsByCreator(c)
}