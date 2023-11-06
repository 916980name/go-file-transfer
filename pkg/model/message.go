package model

import "time"

type Message struct {
	Id          string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Info        string    `json:"info"`
	UserId      string    `json:"userId,omitempty"`
	Create_time time.Time `json:"create_time"`
	Update_time time.Time `json:"update_time"`
}
