package models

import "time"

type Media struct {
	ID        int       `json:"id"`
	ArticleID int       `json:"article_id"`
	FilePath  string    `json:"file_path"`
	FileType  string    `json:"file_type"`
	CreatedAt time.Time `json:"created_at"`
}
