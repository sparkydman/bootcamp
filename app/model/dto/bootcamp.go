package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BootcampRequest struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Phone             string   `json:"phone"`
	Email             string   `json:"email"`
	Address           string   `json:"address"`
	Careers           []string `json:"careers"`
	JobAssistance     bool     `json:"job_assistance"`
	JobGuarantee      bool     `json:"job_guarantee"`
	AcceptPartPayment bool     `json:"accept_part_payment"`
	CreatedBy         string   `json:"created_by"`
	Slug              string   `json:"slug"`
}

func (b BootcampRequest) Validate(validCareers []string) error {
	return validation.ValidateStruct(&b, validation.Field(&b.Name, validation.Required), validation.Field(&b.Description, validation.Required), validation.Field(&b.Phone, validation.Required, validation.Max(20)), validation.Field(&b.Email, validation.Required, is.Email), validation.Field(&b.Address, validation.Required), validation.Field(&b.Careers, validation.Each(validation.In(validCareers))), validation.Field(&b.JobAssistance, validation.Required), validation.Field(&b.JobGuarantee, validation.Required), validation.Field(&b.AcceptPartPayment, validation.Required), validation.Field(&b.CreatedBy, validation.Required, validation.When(!primitive.IsValidObjectID(b.CreatedBy), validation.Required.Error("invalid created_by user id"))))
}

func (b *BootcampRequest) Slugify() {
	b.Slug = slug.Make(b.Name)
}

func ValidCareers() []string {
	return []string{
		"Web Development",
		"Mobile Development",
		"UI/UX",
		"Data Science",
		"Product Management",
		"Project Management",
		"Backend Development",
		"Frontend Development",
		"Graphic Design",
		"Business",
		"Others",
	}
}
