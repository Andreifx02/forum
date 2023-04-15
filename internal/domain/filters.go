package domain

import (
	"time"

	"github.com/google/uuid"
)

type Filters struct {
	DateFrom *time.Time `json:"date_from"`
	DateTo *time.Time `json:"date_to"`
	Topic *string `json:"topic"`
	KeyWords []string `json:"key_words"`
	AuthorID *uuid.UUID `json:"author_id"`
}