package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// The Album struct represends an individual Album from the database.
type Album struct {
	ID      int64   `json:"id"`
	Title   string  `json:"title"`
	Artist  string  `json:"artist"`
	Price   float64 `json:"price"`
	Version int     `json:"-"`
}

// CREATE TABLE IF NOT EXISTS album (
//   id         INT AUTO_INCREMENT NOT NULL,
//   title      VARCHAR(128) NOT NULL,
//   artist     VARCHAR(255) NOT NULL,
//   price      DECIMAL(5,2) NOT NULL,
//   PRIMARY KEY (id)
// )`)

type AlbumModel struct {
	DB *sql.DB
}

func (m AlbumModel) Insert(album *Album) error {
	insertQuery := `
    INSERT INTO album (title, artist, price) 
    VALUES ($1, $2, $3)`

	args := []any{album.Title, album.Artist, album.Price}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := m.DB.ExecContext(ctx, insertQuery, args...)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	album.ID = lastID

	return nil
}

// Returns a slice of movies from the database. Accepts various filter parameters
func (m *AlbumModel) GetAll() ([]*Album, error) {
	query := `
    SELECT id, title, artist, price
    FROM album`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	albums := []*Album{}

	for rows.Next() {
		var album Album

		err := rows.Scan(
			&album.ID,
			&album.Title,
			&album.Artist,
			&album.Price,
		)
		if err != nil {
			return nil, err
		}

		albums = append(albums, &album)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

func (m *AlbumModel) Get(id int64) (*Album, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
    SELECT id, title, artist, price
    FROM album
    WHERE id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var album Album
	args := []any{&album.ID, &album.Title, &album.Artist, &album.Price}

	err := m.DB.QueryRowContext(ctx, query, id).Scan(args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &album, nil
}

func (m *AlbumModel) Update(album *Album) error {
	query := `
    UPDATE movies 
    SET title = ?, artist = ?, price = ?
    WHERE id = ?`

	args := []any{
		album.Title,
		album.Artist,
		album.Price,
		album.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m *AlbumModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
    DELETE FROM album
    WHERE id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
