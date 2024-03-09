package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Character struct {
	ID           int    `json:"ID"`
	CreatedAt    string `json:"CreatedAt"`
	UpdatedAt    string `json:"UpdatedAt"`
	FirstName    string `json:"FirstName"`
	LastName     string `json:"LastName"`
	House        string `json:"House"`
	OriginStatus string `json:"OriginStatus"`
}

type CharacterModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (c CharacterModel) Insert(character *Character) error {
	// Insert a new character item into the database.
	query := `
		INSERT INTO characters (FirstName, LastName, House, OriginStatus) 
		VALUES ($1, $2, $3, $4) 
		RETURNING ID, CreatedAt, UpdatedAt
		`
	args := []interface{}{character.FirstName, character.LastName, character.House, character.OriginStatus}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&character.ID, &character.CreatedAt, &character.UpdatedAt)
}

func (c CharacterModel) Get(id int) (*Character, error) {
	// Retrieve a character item based on its ID.
	query := `	
		SELECT ID, CreatedAt, UpdatedAt, FirstName, LastName, House, OriginStatus
		FROM characters
		WHERE ID = $1
		`
	var character Character
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := c.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&character.ID, &character.CreatedAt, &character.UpdatedAt, &character.FirstName, &character.LastName, &character.House, &character.OriginStatus)
	if err != nil {
		return nil, err
	}
	return &character, nil
}

func (c CharacterModel) Update(character *Character) error {
	// Update a character item in the database.
	query := `
		UPDATE characters
		SET FirstName = $1, LastName = $2, House = $3
		WHERE ID = $4
		RETURNING UpdatedAt
		`
	args := []interface{}{character.FirstName, character.LastName, character.House, character.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&character.UpdatedAt)
}

func (c CharacterModel) Delete(id int) error {
	// Delete a character  item from the database.
	query := `
		DELETE FROM characters
		WHERE ID = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := c.DB.ExecContext(ctx, query, id)
	return err
}
