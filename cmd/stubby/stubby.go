package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	stubby "github.com/shanna/stubby"
	bjf "github.com/shanna/stubby/id/bjf"
	store "github.com/shanna/stubby/store/sqlite"
)

// TODO: flags.
// TODO: errors.

func main() {
	conn, err := store.Open("./stubby.db")
	if err != nil {
		panic(err) // TODO:
	}
	defer conn.Close()

	id, err := bjf.New([]byte("hard coded secret to be replaced"))
	if err != nil {
		panic(err) // TODO:
	}

	shortener := stubby.New(conn, id)

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Mount("/s", shortener.Handler())

	// TODO: Bundle public css/html into binary.
	cwd, _ := os.Getwd()
	web := filepath.Join(cwd, "public")
	router.FileServer("/", http.Dir(web))

	http.ListenAndServe(":3066", router)
}
