package model

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
