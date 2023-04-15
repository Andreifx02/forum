package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	postrgres "github.com/Andreifx02/forum/internal/storage/postgres"
	"github.com/gorilla/mux"
)
	
func (s *Server) GetUserByName(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	username := mux.Vars(r)["username"]

	user, err := s.storage.GetUserByName(ctx, username)

	if err != nil {
		if errors.Is(err, postrgres.UserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)	
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}	
		return 
	}

	userJson, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.Write(userJson)
}
