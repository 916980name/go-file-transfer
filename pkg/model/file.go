package model

import "time"

type FileMeta struct {
	Id        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	Sha       string    `bson:"sha" json:"sha"`
	Size      int64     `bson:"size" json:"size"`
	Location  string    `bson:"location" json:"location"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type UserFile struct {
	Id        string    `bson:"_id,omitempty" json:"_id,omitempty"`
	MetaId    string    `bson:"metaId" json:"metaId"`
	UserId    string    `bson:"userId" json:"userId"`
	Name      string    `bson:"name" json:"name"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}
