package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Racer struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	//foreing key
	TeamId string `json:"teamId"`
}

type RacerModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (r RacerModel) Insert(racer *Racer) error {
	// Insert a new racer item into the database.
	query := `
		INSERT INTO racers (firstName, lastName, teamId) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{racer.FirstName, racer.LastName, racer.TeamId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&racer.Id, &racer.CreatedAt, &racer.UpdatedAt)
}

func (r RacerModel) Get(id int) (*Racer, error) {
	// Retrieve a specific racer item based on its ID.
	query := `
		SELECT id, created_at, updated_at, firstName, lastName, teamId
		FROM racers
		WHERE id = $1
		`
	var racer Racer
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&racer.Id, &racer.CreatedAt, &racer.UpdatedAt, &racer.FirstName, &racer.LastName, &racer.TeamId)
	if err != nil {
		return nil, err
	}
	return &racer, nil
}

func (r RacerModel) Update(racer *Racer) error {
	// Update a specific racer item in the database.
	query := `
		UPDATE racers
		SET firstName = $1, lastName = $2, teamId = $3
		WHERE id = $4
		RETURNING updated_at
		`
	args := []interface{}{racer.FirstName, racer.LastName, racer.TeamId, racer.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return r.DB.QueryRowContext(ctx, query, args...).Scan(&racer.UpdatedAt)
}

func (r RacerModel) Delete(id int) error {
	// Delete a specific racer item from the database.
	query := `
		DELETE FROM racers
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}
