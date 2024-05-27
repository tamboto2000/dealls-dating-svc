package messages

type Profile struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	BirthDate    string `json:"birth_date"`
	Gender       string `json:"gender"`
	RelationNeed string `json:"relation_need"`
	Headline     string `json:"headline,omitempty"`
}
