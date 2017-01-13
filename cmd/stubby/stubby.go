package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	stubby "github.com/techspaceco/techspace-stubby"
)

// TODO: flags.
// TODO: errors.

func main() {
	conn, err := stubby.Open("./stubby.db")
	if err != nil {
		panic(err) // TODO:
	}
	defer conn.Close()

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Mount("/s", stubby.Handler(conn))

	// TODO: Bundle public css/html into binary.
	cwd, _ := os.Getwd()
	web := filepath.Join(cwd, "public")
	router.FileServer("/", http.Dir(web))

	http.ListenAndServe(":3066", router)
}
