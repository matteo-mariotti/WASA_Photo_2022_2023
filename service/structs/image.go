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
	PhotoID        int    `json:"id"`
	Owner          string `json:"photoOwner"`
	Date           string `json:"date"`
	LikesNumber    int    `json:"likes"`
	CommentsNumber int    `json:"comments"`
	LoggedLike     bool   `json:"loggedLike"`
}
