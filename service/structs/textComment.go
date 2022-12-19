package structs

type TextComment struct {
	Text string `json:"text"`
}

type ResponseComment struct {
	CommentID string `json:"id"`
	UserID    string `json:"userID"`
	Text      string `json:"text"`
}
