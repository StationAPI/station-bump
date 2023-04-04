package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stationapi/station-bump/db"
	"github.com/stationapi/station-bump/routes"
)

func main() {
	db, err := db.Connect()

	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	r.Post("/bump/new-bump", func(w http.ResponseWriter, r *http.Request) {
		err := routes.NewBump(w, r, db)

		if err != nil {
			fmt.Println(err)
		}
	})
}
