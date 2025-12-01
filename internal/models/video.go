package models

import (
	"time"
)

type Video struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FilePath    string    `json:"file_path"`
	Thumbnail   string    `json:"thumbnail"`
	Views       int       `json:"views"`
	Likes       int       `json:"likes"`
	Dislikes    int       `json:"dislikes"`
	CreatedAt   time.Time `json:"created_at"`
	Duration    string    `json:"duration"`
	Author      string    `json:"author"`
}

type CreateVideoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
}
