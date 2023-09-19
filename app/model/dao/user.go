package dao

type User struct {
	Base             `bson:"inline"`
	Name             string `json:"name,omitempty" bson:"name"`
	Email            string `json:"email,omitempty" bson:"email"`
	Role             string `json:"role,omitempty" bson:"role"`
	Password         string `json:"password,omitempty" bson:"password"`
	IsConfirmed      bool   `json:"is_confirmed,omitempty" bson:"is_confirmed"`
	TwoFactorEnabled bool   `json:"two_factor_enabled,omitempty" bson:"two_factor_enabled"`
}
