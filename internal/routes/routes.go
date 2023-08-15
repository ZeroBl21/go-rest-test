package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ZeroBl21/go-sql/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Encapsulates the API routes of the application.
func NewRouter(m models.Models) http.Handler {
	handlers := Handlers{Models: m}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", greet)

	r.Route("/api", func(r chi.Router) {
		r.Get("/albums", handlers.listAlbums)
		r.Get("/albums/{id}", handlers.getAlbum)

		r.Post("/albums", handlers.createAlbum)
		r.Put("/albums/{id}", handlers.updateAlbum)

		r.Delete("/albums/{id}", handlers.deleteAlbum)
	})

	return r
}

type Handlers struct {
	models.Models
}

func greet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<main> <h1>Hello World!</h1> <br> <a href='/api/albums'>Go to Albums API</a><main><style> html { filter: invert(1); }</style>")
}
