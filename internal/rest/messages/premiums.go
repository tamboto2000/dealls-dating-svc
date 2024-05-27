package messages

type PremiumFeature struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Types       []string `json:"types"`
}

type PremiumFeaturesResponse []PremiumFeature
