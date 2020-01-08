package cbs

import "testing"

// TODO test functions

func TestLoginCbs(t *testing.T) {
	username := "9802089251"
	password := "k3EfVSamW&W8F^"

	err := LoginCbs(username, password)
	if err != nil {
		t.Errorf("Error while logging in: %s ", err)
	}
}
