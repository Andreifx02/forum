package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Andreifx02/forum/internal/domain"
	"github.com/google/uuid"
)

func (s *Server) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ID uuid.UUID `json:"id"`
		SubID uuid.UUID `json:"sub_id"`
	}
	
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ctx := context.Background()
	err = s.storage.CreateSubscription(ctx, &domain.Subscriptions{
		ID:       req.ID,
		SubID: req.SubID,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
