package structs

type Profile struct {
	Username   string `json:"username"`
	Follower   int    `json:"follower"`
	Following  int    `json:"following"`
	PhotoCount int    `json:"photoNumber"`
}
