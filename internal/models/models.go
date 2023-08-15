package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Container witch can hold and represends all the database Models.
type Models struct {
	Albums AlbumModel
	// Movies      MovieModel
	// Permissions PermissionModel
	// Tokens      TokenModel
	// Users       UserModel
}

// Returns a Models struct containing the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Albums: AlbumModel{DB: db},
		// Movies:      MovieModel{DB: db},
		// Permissions: PermissionModel{DB: db},
		// Tokens:      TokenModel{DB: db},
		// Users:       UserModel{DB: db},
	}
}
