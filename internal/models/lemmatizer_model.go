package models

type Analysis struct {
	Lex string `json:"lex"`
	Gr  string `json:"gr"`
}

type AnalizedWord struct {
	Analysis []Analysis `json:"analysis"`
	Text     string     `json:"text"`
}

type Lemma struct {
	Word, Lemma string
}
