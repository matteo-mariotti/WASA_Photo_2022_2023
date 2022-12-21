package structs

type UserInfo struct {
	User string `json:"identifier"`
}

type Username struct {
	Username string `json:"username"`
}

type Status struct {
	StatusRes bool `json:"status"`
}
