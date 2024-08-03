package models

import "time"

type Comment struct {
	ID        int       `json:"id"`
	ArticleID int       `json:"article_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentLike struct {
	CommentID int `json:"comment_id"`
	UserID    int `json:"user_id"`
}
