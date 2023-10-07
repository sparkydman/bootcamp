package service

import "bootcamp-api/app/repository"

type IBootcampService interface {
}

type bootcamp struct {
	repo repository.IBootcampRepository
}

func NewBootcampService(repo repository.IBootcampRepository) IBootcampService {
	return &bootcamp{repo: repo}
}
