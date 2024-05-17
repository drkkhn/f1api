package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Team struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Name      string `json:"name"`
	Car       string `json:"car"`
}

type TeamModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (t TeamModel) Insert(team *Team) error {
	// Insert a new team item into the database.
	query := `
		INSERT INTO teams (name, car) 
		VALUES ($1, $2) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{team.Name, team.Car}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&team.Id, &team.CreatedAt, &team.UpdatedAt)
}

func (t TeamModel) Get(id int) (*Team, error) {
	// Retrieve a specific team item based on its ID.
	query := `
		SELECT id, created_at, updated_at, name, car
		FROM teams
		WHERE id = $1
		`
	var team Team
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := t.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&team.Id, &team.CreatedAt, &team.UpdatedAt, &team.Name, &team.Car)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (t TeamModel) GetTeamRacers(id int) (*[]Racer, error) {
	query := `
		SELECT id, created_at, updated_at, firstname, lastname, teamid
		FROM racers
		WHERE teamid = $1
		`
	fmt.Println(id)
	var racers []Racer
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := t.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var racer Racer
		if err := rows.Scan(&racer.Id, &racer.CreatedAt, &racer.UpdatedAt, &racer.FirstName, &racer.LastName, &racer.TeamId); err != nil {
			return nil, err
		}
		racers = append(racers, racer)
	}
	fmt.Println(racers)
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &racers, nil
}

func (t TeamModel) Update(team *Team) error {
	// Update a specific team item in the database.
	query := `
		UPDATE teams
		SET name = $1, car = $2
		WHERE id = $3
		RETURNING updated_at
		`
	args := []interface{}{team.Name, team.Car, team.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&team.UpdatedAt)
}

func (t TeamModel) Delete(id int) error {
	// Delete a specific team item from the database.
	query := `
		DELETE FROM teams
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := t.DB.ExecContext(ctx, query, id)
	return err
}
