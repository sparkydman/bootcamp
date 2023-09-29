package dto

import (
	"bootcamp-api/app/model/dao"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserRequest struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type CreateUserRequest struct {
	UserRequest
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginUserResponse struct {
	dao.User
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (c CreateUserRequest) Vaildate() error {
	return validation.ValidateStruct(&c, validation.Field(&c.Name, validation.Required), validation.Field(&c.Email, validation.Required, is.Email), validation.Field(&c.Role, validation.Required, validation.In("user", "publisher", "admin")), validation.Field(&c.Password, validation.Required, validation.Length(8, 30)), validation.Field(&c.ConfirmPassword, validation.When(c.Email != c.ConfirmPassword, validation.Required.Error("Password don't match"))))
}
func (c UserRequest) Vaildate() error {
	return validation.ValidateStruct(&c, validation.Field(&c.Name, validation.Required), validation.Field(&c.Role, validation.Required, validation.In("user", "publisher", "admin")))
}
func (l LoginUserRequest) Vaildate() error {
	return validation.ValidateStruct(&l, validation.Field(&l.Email, validation.Required, is.Email), validation.Field(&l.Password, validation.Required, validation.Length(8, 30)))
}
