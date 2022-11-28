package structs

type Like struct {
	UserID string `json:"userID"`
}

type Comment struct {
	CommentID string `json:"commentID"`
	UserID    string `json:"userID"`
	Text      string `json:"text"`
}

type Photo struct {
	PhotoID  string    `json:"photoID"`
	Likes    []Like    `json:"likes"`
	Comments []Comment `json:"comments"`
}
