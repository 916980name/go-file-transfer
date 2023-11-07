package model

import "time"

type Message struct {
	Id        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Info      string    `json:"info"`
	UserId    string    `json:"userId,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
