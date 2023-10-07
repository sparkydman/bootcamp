package dependencies

import (
	"bootcamp-api/app/controller"
	"bootcamp-api/app/repository"
	"bootcamp-api/app/service"
	"bootcamp-api/config"
)

type Initialize struct {
	Db *config.Db
	UserRepo repository.IUserRepository
	UserSvc  service.IUserService
	UserCtrl controller.IUserController
}

func NewInitialize(db *config.Db, userRepo repository.IUserRepository, userSvc service.IUserService, userCtrl controller.IUserController) *Initialize {
	return &Initialize{Db: db, UserRepo: userRepo, UserSvc: userSvc, UserCtrl: userCtrl}
}
