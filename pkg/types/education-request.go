package types

import "time"

type EducationRequest struct {
	ID            uint      `json:"id"`
	School        string    `json:"school"`
	Degree        string    `json:"degree"`
	Fieldofstudy  string    `json:"fieldofstudy"`
	From          string    `json:"from"`
	ConvertedFrom time.Time `json:"-"`
	To            string    `json:"to"`
	ConvertedTo   time.Time `json:"-"`
	Current       bool      `json:"current"`
	Description   string    `json:"description"`
}
