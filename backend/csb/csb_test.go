package csb

import "testing"

// TODO test functions

/* func TestHandleGreet(t *testing.T) {
	is := is.New(t)
	srv := newServer()
	db, cleanup := connectTestDatabase()
	defer cleanup()
	srv.db = db
	r := httptest.NewRequest("GET", "/greet", nil)
	w := httptest.NewRecorder()
	srv.handleGreet(w, r)
	is.Equal(w.Code, http.StatusOK)
}
*/

func TestLoginCbs(t *testing.T) {
	username := "invaliduser"
	password := "invalidpassword"

	err := LoginCbs(username, password)
	if err == nil {
		t.Error("expecting error when using a wrong login")
	}
}
