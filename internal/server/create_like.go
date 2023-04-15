package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Andreifx02/forum/internal/domain"
	"github.com/google/uuid"
)
	
func (s *Server) CreateLike(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserID uuid.UUID `json:"user_id"`
		PostID uuid.UUID `json:"post_id"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	err = s.storage.CreateLike(ctx, &domain.Like{
		UserID:    req.UserID,
		PostID:    req.PostID,
		
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
