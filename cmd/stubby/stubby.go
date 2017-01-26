package main

//go:generate go-bindata -prefix "public/" -pkg main -o public.go

import (
	"crypto/rand"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	stubby "github.com/shanna/stubby"
	bjf "github.com/shanna/stubby/id/bjf"
	store "github.com/shanna/stubby/store/sqlite"
)

// TODO: errors.

var config struct {
	secret   []byte
	listen   string
	database string
}

func main() {
	secret := os.Getenv("STUBBY_SECRET")
	if secret != "" {
		config.secret = []byte(secret)
	} else {
		fmt.Println("stubby warning: STUBBY_SECRET environment variable not set. Using random secret.")
		config.secret = make([]byte, 32)
		rand.Read(config.secret)
	}

	flag.StringVar(&config.listen, "listen", "localhost:3030", "listen address")
	flag.StringVar(&config.database, "database", "./stubby.db", "sqlite database path")
	flag.Parse()

	conn, err := store.Open(config.database)
	if err != nil {
		panic(err) // TODO:
	}
	defer conn.Close()

	id, err := bjf.New(config.secret)
	if err != nil {
		panic(err) // TODO:
	}

	shortener, err := stubby.New(conn, id)
	if err != nil {
		panic(err) // TODO:
	}

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Mount("/s", stubby.Handler(shortener, config.secret))

	// TODO: Bundle public css/html into binary.
	cwd, _ := os.Getwd()
	web := filepath.Join(cwd, "public")
	router.FileServer("/", http.Dir(web))

	http.ListenAndServe(":3066", router)
}
