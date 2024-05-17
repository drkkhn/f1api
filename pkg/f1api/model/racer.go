package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/drkkhn/f1api/pkg/f1api/validator"
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

func ValidateRacer(v *validator.Validator, racer *Racer) {
	v.Check(racer.FirstName != "", "firstname", "must be provided")
	v.Check(len(racer.FirstName) <= 500, "firstname", "must not be more than 500 bytes long")
}

func ValidateTeam(v *validator.Validator, team *Team) {
	v.Check(team.Name != "", "name", "must be provided")
	v.Check(len(team.Name) <= 500, "name", "must not be more than 500 bytes long")
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

func (r RacerModel) GetAll(firstname string, teamid string, filters Filters) ([]*Racer, Metadata, error) {
	if teamid == "" {
		teamid = "0"
	}
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, updated_at, firstname, lastname, teamid
		FROM racers
		WHERE (LOWER(firstname) = LOWER($1) OR $1 = '')
		AND (teamid >= $2 OR $2 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{firstname, teamid, filters.limit(), filters.offset()}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	racers := []*Racer{}

	for rows.Next() {
		var racer Racer

		err := rows.Scan(
			&totalRecords,
			&racer.Id,
			&racer.CreatedAt,
			&racer.UpdatedAt,
			&racer.FirstName,
			&racer.LastName,
			&racer.TeamId,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		racers = append(racers, &racer)
	}

	if err := rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return racers, metadata, nil // Return racers and nil error
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
