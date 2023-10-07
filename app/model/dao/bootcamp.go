package dao

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bootcamp struct {
	Name              string             `json:"name,omitempty" bson:"name"`
	Slug              string             `json:"slug,omitempty" bson:"slug"`
	Description       string             `json:"description,omitempty" bson:"description"`
	Website           string             `json:"website,omitempty" bson:"website"`
	Phone             string             `json:"phone,omitempty" bson:"phone"`
	Email             string             `json:"email,omitempty" bson:"email"`
	Address           string             `json:"address,omitempty" bson:"address"`
	Careers           []string           `json:"careers,omitempty" bson:"careers"`
	AverageRating     int                `json:"average_rating,omitempty" bson:"average_rating"`
	JobAssistance     bool               `json:"job_assistance,omitempty" bson:"job_assistance"`
	JobGuarantee      bool               `json:"job_guarantee,omitempty" bson:"job_guarantee"`
	AcceptPartPayment bool               `json:"accept_part_payment,omitempty" bson:"accept_part_payment"`
	CreatedBy         primitive.ObjectID `json:"created_by,omitempty" bson:"created_by"`
	Base
}
