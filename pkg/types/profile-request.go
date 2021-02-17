package types

import "time"

type Experience struct {
	Title         string    `json:"title"`
	Company       string    `json:"company"`
	Location      string    `json:"location"`
	From          string    `json:"from"`
	ConvertedFrom time.Time `json:"-"`
	To            string    `json:"to"`
	ConvertedTo   time.Time `json:"-"`
	Current       bool      `json:"current"`
	Description   string    `json:"description"`
}

type Social struct {
	Youtube   string `json:"youtube"`
	Twitter   string `json:"twitter"`
	Facebook  string `json:"facebook"`
	Linkedin  string `json:"linkedin"`
	Instagram string `json:"instagram"`
}

type ProfileRequest struct {
	Company        string       `json:"company"`
	Website        string       `json:"website"`
	Location       string       `json:"location"`
	Status         string       `json:"status"`
	Skills         string       `json:"skills"`
	SkillList      []string     `json:"-"`
	Bio            string       `json:"bio"`
	Githubusername string       `json:"githubusername"`
	Social         Social       `json:"social"`
	Experience     []Experience `json:"experience"`
	UserID         int
}
