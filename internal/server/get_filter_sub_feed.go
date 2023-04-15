package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Andreifx02/forum/internal/domain"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) GetFilterSubFeed(w http.ResponseWriter, r *http.Request) {
	var request domain.Filters
	idStr := mux.Vars(r)["id"]
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	
	if err != nil {
		http.Error(w, fmt.Errorf("Invalid ID: %w", err).Error(), http.StatusUnprocessableEntity)
		return
	}

	ctx := context.Background()
	posts, err := s.storage.GetFilterSubFeed(ctx, id, &request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postJson, _ := json.Marshal(posts)
	w.Header().Set("Content-Type", "application/json")
	w.Write(postJson)
	
}