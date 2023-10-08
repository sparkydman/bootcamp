package dto

import (
	"errors"
	"fmt"
	"strings"

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
	return validation.ValidateStruct(&b, validation.Field(&b.Name, validation.Required), validation.Field(&b.Description, validation.Required), validation.Field(&b.Phone, validation.Required, validation.Length(12, 20)), validation.Field(&b.Email, validation.Required, is.Email), validation.Field(&b.Address, validation.Required), validation.Field(&b.Careers, validation.Required, validation.Each(validation.By(customInRule(ValidCareers())))), validation.Field(&b.CreatedBy, validation.Required, validation.When(!primitive.IsValidObjectID(b.CreatedBy), validation.Required.Error("invalid created_by user id"))))
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

func customInRule(data []string) validation.RuleFunc {
	mapdata := make(map[string]string, len(data))
	for _, v := range data {
		newValue := strings.ToLower(v)
		mapdata[newValue] = v
	}

	return func(value interface{}) error {
		s, ok := value.(string)
		if !ok {
			return errors.New("invalid string value")
		}
		if _, ok := mapdata[strings.ToLower(s)]; !ok {
			return errors.New(fmt.Sprintf("%s: is not a valid option", s))
		}
		return nil
	}
}
