package models

type UserAuth struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
