package controller

import (
	"bootcamp-api/app/service"

	"github.com/gin-gonic/gin"
)

type IUserController interface {
	CreateUser(c *gin.Context)
	GetUserById(c *gin.Context)
	LoginUser(c *gin.Context)
	GetUsers(c *gin.Context)
	UpdateUser(g *gin.Context)
	DeleteUser(g *gin.Context)
}

type userController struct {
	svc service.IUserService
}

func NewUserController(svc service.IUserService) IUserController {
	return &userController{svc: svc}
}

func (c userController) CreateUser(g *gin.Context) {
	c.svc.CreateUser(g)
}

func (c userController) GetUserById(g *gin.Context) {
	c.svc.GetUserById(g)
}
func (c userController) LoginUser(g *gin.Context) {
	c.svc.LoginUser(g)
}
func (c userController) GetUsers(g *gin.Context) {
	c.svc.GetUsers(g)
}
func (c userController) UpdateUser(g *gin.Context) {
	c.svc.UpdateUser(g)
}
func (c userController) DeleteUser(g *gin.Context) {
	c.svc.DeleteUser(g)
}
