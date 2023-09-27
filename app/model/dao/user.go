package dao

type User struct {
	Base             `bson:"inline"`
	Name             string `json:"name" bson:"name"`
	Email            string `json:"email" bson:"email"`
	Role             string `json:"role" bson:"role"`
	Password         string `json:"password,omitempty" bson:"password"`
	IsConfirmed      bool   `json:"is_confirmed" bson:"is_confirmed"`
	TwoFactorEnabled bool   `json:"two_factor_enabled" bson:"two_factor_enabled"`
}
