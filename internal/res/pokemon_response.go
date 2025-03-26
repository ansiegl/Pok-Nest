package res

// response structure
type PokemonResponse struct {
	PokemonID  string `json:"pokemon_id"`
	Name       string `json:"name"`
	Type1      string `json:"type_1"`
	Type2      string `json:"type_2,omitempty"`
	Generation int    `json:"generation"`
	Legendary  bool   `json:"legendary"`
}

type PaginationMetadata struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type APIResponse struct {
	Data       []PokemonResponse  `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}
