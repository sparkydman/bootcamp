package controller

import "bootcamp-api/app/service"

type IBootcampController interface {
}

type bootcampController struct {
	svc service.IBootcampService
}

func NewBootcampController(svc service.IBootcampService) IBootcampController {
	return &bootcampController{svc: svc}
}
