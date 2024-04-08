package table

type Column struct {
	Name    string `json:"name"`
	Default string `json:"default"`
	Type    string `json:"type"`
}
