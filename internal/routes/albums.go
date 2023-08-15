package routes

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ZeroBl21/go-sql/internal/models"
	"github.com/ZeroBl21/go-sql/internal/utils"
)

func (h *Handlers) listAlbums(w http.ResponseWriter, r *http.Request) {
	albums, err := h.Models.Albums.GetAll()
	if err != nil {
		http.Error(w, "Error fetching albums", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, 200, utils.Envelope{"albums": albums}, nil)
}

func (h *Handlers) getAlbum(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParams(r)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	album, err := h.Models.Albums.Get(id)
	if err != nil {
		http.Error(w, "Error fetching albums", http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"album": album}, nil)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
	}
}

func (h *Handlers) createAlbum(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string  `json:"title"`
		Artist string  `json:"artist"`
		Price  float64 `json:"price"`
	}

	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	album := &models.Album{
		Title:  input.Title,
		Artist: input.Artist,
		Price:  input.Price,
	}

	err = h.Models.Albums.Insert(album)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", album.ID))

	err = utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"album": album}, headers)
	if err != nil {
		log.Println("JSON", err)
		http.Error(w, "Error", http.StatusInternalServerError)
	}
}

func (h *Handlers) updateAlbum(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParams(r)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	album, err := h.Models.Albums.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			http.Error(w, "Not Found", http.StatusNotFound)
		default:
			http.Error(w, "Error", http.StatusInternalServerError)
		}
		return
	}

	var input struct {
		Title  string  `json:"title"`
		Artist string  `json:"artist"`
		Price  float64 `json:"price"`
	}

	err = utils.ReadJSON(w, r, &input)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if input.Title != "" {
		album.Title = input.Title
	}
	if input.Artist != "" {
		album.Artist = input.Artist
	}
	if input.Price != 0.0 {
		album.Price = input.Price
	}

	err = h.Models.Albums.Update(album)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"album": album}, nil)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
	}
}

func (h *Handlers) deleteAlbum(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParams(r)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	err = h.Models.Albums.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			http.Error(w, "Not Found", http.StatusNotFound)
		default:
			http.Error(w, "internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "album successfully deleted"}, nil)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
	}
}
