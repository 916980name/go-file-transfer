package model

type UserInfo struct {
	Id        string `bson:"_id,omitempty" json:"_id,omitempty"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	PostCount int64  `json:"postCount"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type BizUserInfo struct {
	Username   string   `json:"username"`
	Privileges []string `json:"privileges"`
}
