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
	GetLoggedInUser(c *gin.Context)
	GetToken(c *gin.Context)
}

type UserController struct {
	svc service.IUserService
}

func NewUserController(svc service.IUserService) *UserController {
	return &UserController{svc: svc}
}

func (c UserController) CreateUser(g *gin.Context) {
	c.svc.CreateUser(g)
}

func (c UserController) GetUserById(g *gin.Context) {
	c.svc.GetUserById(g)
}
func (c UserController) LoginUser(g *gin.Context) {
	c.svc.LoginUser(g)
}
func (c UserController) GetUsers(g *gin.Context) {
	c.svc.GetUsers(g)
}
func (c UserController) UpdateUser(g *gin.Context) {
	c.svc.UpdateUser(g)
}
func (c UserController) DeleteUser(g *gin.Context) {
	c.svc.DeleteUser(g)
}
func (c UserController) GetLoggedInUser(g *gin.Context) {
	c.svc.GetLoggedInUser(g)
}
func (c UserController) GetToken(g *gin.Context) {
	c.svc.GetToken(g)
}
