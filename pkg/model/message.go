package model

import "time"

type Message struct {
	Id        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Info      string    `bson:"info" json:"info"`
	UserId    string    `bson:"userId" json:"userId"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type ShareMessage struct {
	Id        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	MessageId string    `bson:"messageId" json:"messageId"`
	UserId    string    `bson:"userId" json:"userId"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
