package models

import "time"

type Article struct {
	ID        int       `json:"id"`
	AuthorID  int       `json:"author"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ArticleCategory struct {
	ArticleID int    `json:"article_id"`
	Category  string `json:"category_id"`
}

type ArticleTag struct {
	ArticleID int `json:"article_id"`
	TagID     int `json:"tag_id"`
}

type ArticleLike struct {
	ArticleID int `json:"article_id"`
	UserID    int `json:"user_id"`
}
