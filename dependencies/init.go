package dependencies

import (
	"bootcamp-api/app/controller"
	"bootcamp-api/app/repository"
	"bootcamp-api/app/service"
	"bootcamp-api/config"
)

type Initialize struct {
	//initialize db
	Db *config.Db
	//initialize user
	UserRepo repository.IUserRepository
	UserSvc  service.IUserService
	UserCtrl controller.IUserController
	//initialize bootcamp
	BootcampRepo repository.IBootcampRepository
	BootcampSvc  service.IBootcampService
	BootcampCtrl controller.IBootcampController
}

func NewInitialize(db *config.Db, userRepo repository.IUserRepository, userSvc service.IUserService, userCtrl controller.IUserController, bootcampRepo repository.IBootcampRepository, bootcampSvc service.IBootcampService, bootcampCtrl controller.IBootcampController) *Initialize {
	return &Initialize{Db: db, UserRepo: userRepo, UserSvc: userSvc, UserCtrl: userCtrl, BootcampRepo: bootcampRepo, BootcampSvc: bootcampSvc, BootcampCtrl: bootcampCtrl}
}
