package structs

type profile struct {
	Username    string `json:"username"`
	Follower    int    `json:"follower"`
	Following   int    `json:"following"`
	PhotoNumber int    `json:"photo"`
}
