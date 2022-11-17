package structs

type Fountain struct {
	ID        int     `json:"id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
}

func (fountain Fountain) IsValidFountain() bool {
	//Example of validation
	return fountain.ID > 0
}
