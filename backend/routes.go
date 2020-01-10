package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/vifraa/opencsb/csb"
)

func (s *server) routes() {
	s.router = chi.NewMux()

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)

	s.router.Get("/", s.handleIndex())

	// TODO Instead of taking the login info each time in post request, have a login endpoint instead.
	s.router.Post("/doors/open", s.handleDoorOpen())
	s.router.Post("/doors", s.handleDoorGetAll())
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, "Hello World", http.StatusOK)
	}
}

func (s *server) handleDoorOpen() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		doorReq := &DoorOpenRequest{}
		err := s.decode(w, r, doorReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = csb.LoginCbs(doorReq.Username, doorReq.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = csb.LoginAptusPort()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = csb.OpenDoor(doorReq.DoorID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.respond(w, r, "Door opened", http.StatusOK)
	}
}

// TODO Alot of code duplication between these two methods, need to refactor.
func (s *server) handleDoorGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		doorReq := &DoorOpenRequest{}
		err := s.decode(w, r, doorReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = csb.LoginCbs(doorReq.Username, doorReq.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = csb.LoginAptusPort()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		doors, err := csb.FetchDoorIDs()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		s.respond(w, r, doors, http.StatusOK)
	}
}

type DoorOpenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DoorID   string `json:"doorId"`
}
