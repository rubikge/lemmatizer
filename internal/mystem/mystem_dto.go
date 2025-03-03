package mystem

type AnalizedWord struct {
	Analysis []Analysis `json:"analysis"`
	Text     string     `json:"text"`
}

type Analysis struct {
	Lex string `json:"lex"`
	Gr  string `json:"gr"`
}
