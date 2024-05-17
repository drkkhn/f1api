package filler

import (
	model "github.com/drkkhn/f1api/pkg/f1api/model"
)

func PopulateDatabase(models model.Models) error {
	for _, racer := range racers {
		models.Racers.Insert(&racer)
	}
	// TODO: Implement restaurants pupulation
	// TODO: Implement the relationship between restaurants and menus
	return nil
}

var racers = []model.Racer{
	{FirstName: "John", LastName: "Doe", TeamId: "Team1"},
	{FirstName: "Jane", LastName: "Smith", TeamId: "Team2"},
	{FirstName: "Alice", LastName: "Johnson", TeamId: "Team3"},
	{FirstName: "Bob", LastName: "Williams", TeamId: "Team1"},
	{FirstName: "Emily", LastName: "Brown", TeamId: "Team2"},
	{FirstName: "Michael", LastName: "Wilson", TeamId: "Team3"},
	{FirstName: "Emma", LastName: "Jones", TeamId: "Team1"},
	{FirstName: "Daniel", LastName: "Martinez", TeamId: "Team2"},
	{FirstName: "Olivia", LastName: "Taylor", TeamId: "Team3"},
	{FirstName: "William", LastName: "Anderson", TeamId: "Team1"},
}
