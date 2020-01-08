package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
	"github.com/vifraa/opencbs/cbs"
)

func (s *server) routes() {
	s.router = chi.NewMux()

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)

	s.router.Get("/", s.handleIndex())

	s.router.Post("/door", s.handleDoorOpen())

	//exerciseStore := SqlExerciseStore{
	//	db: s.db,
	//}
	//s.router.Get("/exercises/{id}", s.handleExerciseGet(exerciseStore))
	//s.router.Post("/exercises", s.handleExerciseCreate(exerciseStore))
}

func (s *server) handleIndex() http.HandlerFunc {
	//preparations here, can insert arguments
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, "Hello World", http.StatusOK)
	}
}

func (s *server) handleDoorOpen() http.HandlerFunc {
	type DoorOpenRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		DoorID   string `json:"doorId"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		doorReq := &DoorOpenRequest{}
		err := s.decode(w, r, doorReq)
		if err != nil {
			s.respond(w, r, errors.Wrap(err, "invalid json"), http.StatusBadRequest)
		}

		err = cbs.LoginCbs(doorReq.Username, doorReq.Password)
		if err != nil {
			s.respond(w, r, err, http.StatusBadRequest)
		}

		err = cbs.LoginAptusPort()
		if err != nil {
			s.respond(w, r, err, http.StatusBadRequest)
		}

		err = cbs.OpenDoor(doorReq.DoorID)
		if err != nil {
			s.respond(w, r, err, http.StatusBadRequest)
		}

		s.respond(w, r, "Door opened", http.StatusBadRequest)
	}
}

//handleTaskCreate
//handleTaskDone
//handleTaskGet
//...
