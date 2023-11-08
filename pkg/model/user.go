package model

import "time"

type UserInfo struct {
	Id        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Username  string    `bson:"username" json:"username"`
	Password  string    `bson:"password" json:"password"`
	Email     string    `bson:"email"  json:"email"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type BizUserInfo struct {
	Username   string   `json:"username"`
	Privileges []string `json:"privileges"`
}
