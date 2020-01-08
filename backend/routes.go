package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *server) routes() {
	s.router = chi.NewMux()

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)

	s.router.Get("/", s.handleIndex())

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

//handleTaskCreate
//handleTaskDone
//handleTaskGet
//...
