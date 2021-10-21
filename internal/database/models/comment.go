package models

import (
	"time"
)

type Comment struct {
	ID        string `gorm:"type:uuid;primary_key"`
	Slug      string
	Body      string
	Author    string
	CreatedAt time.Time
}

func NewComment() *Comment {
	return &Comment{}
}
