package server

import (
	"fmt"
	"net/http"

	"github.com/Andreifx02/forum/internal/config"
	postrgres "github.com/Andreifx02/forum/internal/storage/postgres"
	"github.com/gorilla/mux"
)

type Server struct {
	storage *postrgres.Storage
}

func New(storage *postrgres.Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func (s *Server) StartListen(cfg *config.Config) {
	router := mux.NewRouter()

	router.HandleFunc("/user/create", s.CreateUser).Methods("POST")
	router.HandleFunc("/post/create", s.CreatePost).Methods("POST")
	router.HandleFunc("/sub/create", s.CreateSubscription).Methods("POST")
	router.HandleFunc("/like/create", s.CreateLike).Methods("POST")
	router.HandleFunc("/user/{id}/subfeed", s.GetSubFeed).Methods("GET")
	router.HandleFunc("/user/{id}/interesting", s.GetInteresting).Methods("GET")

	

	http.Handle("/", router)

	url := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	err := http.ListenAndServe(url, nil)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
}
