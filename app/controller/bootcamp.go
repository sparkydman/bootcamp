package controller

import (
	"bootcamp-api/app/service"

	"github.com/gin-gonic/gin"
)

type IBootcampController interface {
	AddBootcamp(c *gin.Context)
}

type bootcampController struct {
	svc service.IBootcampService
}

func NewBootcampController(svc service.IBootcampService) IBootcampController {
	return &bootcampController{svc: svc}
}
func (b bootcampController) AddBootcamp(c *gin.Context){
	b.svc.AddBootcamp(c)
}