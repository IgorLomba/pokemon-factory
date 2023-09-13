package models

// Pokemon model
// swagger:response Pokemon
type Pokemon struct {
	Name         string   `json:"name"`
	Capabilities []string `json:"capabilities"`
}
