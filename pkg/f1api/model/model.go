package model

import (
	"database/sql"
	"log"
	"os"
)

type Models struct {
	Racers RacerModel
	Teams  TeamModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Racers: RacerModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Teams: TeamModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
