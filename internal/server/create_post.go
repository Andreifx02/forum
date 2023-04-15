package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Andreifx02/forum/internal/domain"
	"github.com/google/uuid"
)
	
func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	type request struct {
		AuthorID uuid.UUID `json:"author_id"`
		Topic    string    `json:"topic"`
		Text     string    `json:"text"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = s.storage.CreatePost(ctx, &domain.Post{
		ID:       uuid.New(),
		AuthorID: req.AuthorID,
		Topic:    req.Topic,
		Text:     req.Text,
		Date:     time.Now(),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
