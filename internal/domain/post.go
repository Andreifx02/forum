package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct{
	ID uuid.UUID `json:"id"`
	AuthorID uuid.UUID `json:"author_id"`
	Topic string `json:"topic"`
	Text string `json:"text"`
	Date time.Time `json:"date"`
}