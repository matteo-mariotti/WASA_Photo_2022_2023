package structs

type Profile struct {
	Username    string `json:"username"`
	Follower    int    `json:"follower"`
	Following   int    `json:"following"`
	PhotoNumber int    `json:"photoNumber"`
	Photo       []int  `json:"photoIDs"`
}
