package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Andreifx02/forum/internal/domain"
	"github.com/google/uuid"
)

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Nickname string `json:"nickname"`
	}

	var req request
	json.NewDecoder(r.Body).Decode(&req)

	ctx := context.Background()
	err := s.storage.CreateUser(ctx, &domain.User{
		ID:       uuid.New(),
		Nickname: req.Nickname,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
