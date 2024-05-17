package model

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

type Models struct {
	Racers      RacerModel
	Teams       TeamModel
	Permissions PermissionModel
	Tokens      TokenModel
	Members     MemberModel
}

var (
	// ErrRecordNotFound is returned when a record doesn't exist in database.
	ErrRecordNotFound = errors.New("record not found")

	// ErrEditConflict is returned when a there is a data race, and we have an edit conflict.
	ErrEditConflict = errors.New("edit conflict")
)

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Racers: RacerModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Permissions: PermissionModel{
			DB: db,
		},
		Teams: TeamModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Tokens: TokenModel{
			DB: db,
		},
		Members: MemberModel{
			DB: db,
		},
	}
}
