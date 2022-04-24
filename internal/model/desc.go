package model

type ChangeDesc struct {
	Id             int    `json:"id"`
	GroupId        int    `json:"groupid"`
	GroupName      string `json:"groupname"`
	NextDescChange string `json:"nextdescchange"`
	CreatedAt      string `json:"created_at"`
}

type Desc struct {
	Id        int
	Text      string
	CreatedAt string
}
