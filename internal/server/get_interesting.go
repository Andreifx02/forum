package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)
	
func (s *Server) GetInteresting(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	id, err := uuid.Parse(idStr)
	
	if err != nil {
		http.Error(w, fmt.Errorf("Invalid ID: %w", err).Error(), http.StatusUnprocessableEntity)
	}

	ctx := context.Background()
	posts, err := s.storage.GetInteresting(ctx, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	postJson, _ := json.Marshal(posts)
	w.Header().Set("Content-Type", "application/json")
	w.Write(postJson)
}