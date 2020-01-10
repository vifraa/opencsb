package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vifraa/opencbs/cbs"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	srv := newServer()

	s := &http.Server{
		Addr:           "0.0.0.0:8080",
		Handler:        srv,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Server running on: " + s.Addr)

	return s.ListenAndServe()
}

func testOpenDoor() {
	err := cbs.LoginCbs("9802089251", "k3EfVSamW&W8F^")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("logged in")

	err = cbs.LoginAptusPort()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("aptus port logged in")
	//	err = cbs.OpenDoor("123640")
	//	if err != nil {
	//		log.Fatal(err)
	//	}

	ids, err := cbs.FetchDoorIDs()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found door ids: %s", ids)

}
