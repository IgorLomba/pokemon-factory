package models

// RequestStatus model
// swagger:response RequestStatus
type RequestStatus struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Completed bool   `json:"completed"`
}

// PokemonQueryParameters request
// swagger:parameters PokemonQueryParameters
type PokemonQueryParameters struct {
	// ID of the request
	// in: query
	// required: true
	ID string `json:"id"`
}

// PokemonGenerateQueryParameters request
// swagger:parameters PokemonGenerateQueryParameters
type PokemonGenerateQueryParameters struct {
	// Amount to generate
	// in: query
	// required: true
	Amount string `json:"amount"`
}
