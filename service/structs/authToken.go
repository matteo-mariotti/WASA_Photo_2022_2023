package structs

type AuthToken struct {
	Token string `json:"token"`
}

// Adding the stringer interface to the AuthToken struct
func (a AuthToken) String() string {
	return a.Token
}
