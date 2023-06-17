package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var msg = make(chan interface{})

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Content-Type", "text/event-stream")

		for v := range msg {
			fmt.Fprintf(w, "data: %v \n\n", v)
			w.(http.Flusher).Flush()
			time.Sleep(time.Second)
		}
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {

		bd, err := io.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		str := string(bd)

		msg <- strings.TrimSpace(str)

		fmt.Println(strings.TrimSpace(str))

		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", r)
}
