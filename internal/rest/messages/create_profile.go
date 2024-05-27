package messages

type Hobby struct {
	HobbyName   string `json:"hobby_name"`
	Description string `json:"description,omitempty"`
}

type Interest struct {
	InterestedIn string `json:"interested_in"`
}

type CreateProfileRequest struct {
	BirthDate         string     `json:"birth_date"`
	Gender            string     `json:"gender"`
	RelationNeed      string     `json:"relation_need"`
	LastEducation     string     `json:"last_education"`
	LastEducationInst string     `json:"last_education_inst"`
	Headline          string     `json:"headline"`
	Description       string     `json:"description"`
	Hobbies           []Hobby    `json:"hobbies"`
	Interests         []Interest `json:"interests"`
}
