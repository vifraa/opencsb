package main

import (
	"log"

	"github.com/vifraa/opencbs/cbs"
)

func main() {
	err := cbs.LoginCbs("9802089251", "k3EfVSamW&W8F^")
	if err != nil {
		log.Fatal(err)
	}
	err = cbs.LoginAptusPort()
	if err != nil {
		log.Fatal(err)
	}

	err = cbs.OpenDoor("123640")
	if err != nil {
		log.Fatal(err)
	}
}
