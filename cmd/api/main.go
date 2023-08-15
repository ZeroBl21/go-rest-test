package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ZeroBl21/go-sql/internal/models"
	"github.com/ZeroBl21/go-sql/internal/routes"
	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

// Holds the dependencies of the config and routes.
type Application struct {
	// config config
	routes http.Handler
}

func (app *Application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8080),
		Handler:      app.routes,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Println("Running on http://localhost:8080")
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func main() {
	db, err := dbInit()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	models := models.NewModels(db)

	app := &Application{
		routes: routes.NewRouter(models),
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}

const dbUrl = "file:./file.db"

func dbInit() (*sql.DB, error) {
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Query(`
    CREATE TABLE IF NOT EXISTS album (
      id         INTEGER PRIMARY KEY AUTOINCREMENT,
      title      VARCHAR(128) NOT NULL,
      artist     VARCHAR(255) NOT NULL,
      price      DECIMAL(5,2) NOT NULL
    );`)

	if err != nil {
		return nil, err
	}

	return db, err
}
