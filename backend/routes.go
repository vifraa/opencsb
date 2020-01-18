package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/securecookie"
	"github.com/vifraa/opencsb/csb"
	"github.com/vifraa/opencsb/user"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

const sessionName = "opencsb-session"

func (s *server) routes() {
	s.router = chi.NewMux()

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Recoverer)

	s.router.Get("/", s.handleIndex())

	s.router.Post("/login", s.handleLogin())

	// TODO Instead of taking the login info each time in post request, have a login endpoint instead.
	s.router.Post("/doors/open", s.handleDoorOpen())
	s.router.Post("/doors", s.handleDoorGetAll())
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, "Hello World", http.StatusOK)
	}
}

func (s *server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")

		err := csb.LoginCbs(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO encrypt password!
		u := user.User{
			Username:  username,
			Password:  password,
			LastLogin: time.Now(),
		}
		s.userRepository.Update(u)

		setSession(username, w)
	}
}

func setSession(username string, w http.ResponseWriter) {
	v := map[string]string{
		"username": username,
	}

	encoded, err := cookieHandler.Encode(sessionName, v)
	if err != nil {
		return
	}

	cookie := &http.Cookie{
		Name:  sessionName,
		Value: encoded,
		Path:  "/",
	}

	http.SetCookie(w, cookie)
}

func getUsernameFromCookie(r *http.Request) string {
	cookie, err := r.Cookie(sessionName)
	if err != nil {
		// TODO handle error
		return ""
	}

	cv := make(map[string]string)
	err = cookieHandler.Decode(sessionName, cookie.Value, &cv)
	if err != nil {
		// TODO handle error
		return ""
	}

	return cv["username"]
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
