package model

import "time"

type UserInfo struct {
	Id        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type BizUserInfo struct {
	Username   string   `json:"username"`
	Privileges []string `json:"privileges"`
}
