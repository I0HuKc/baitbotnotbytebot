package model

type Performance struct {
	Id        int    `json:"id"`
	GroupId   int    `json:"groupid"`
	GroupName string `json:"groupname"`
	NextJoke  string `json:"nextjoke"`
	CreatedAt string `json:"created_at"`
}

type Joke struct {
	Id       string
	Target   string
	Setup    string
	Delivery string
	Lang     JLang
}

type JLang struct {
	Source    string
	Target    string
	Translate bool
}
