package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Character struct {
	ID           int    `json:"id"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
	FirstName    string `json:"fisrtName"`
	LastName     string `json:"lastName"`
	House        string `json:"house"`
	OriginStatus string `json:"originStatus"`
}

type CharacterModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (c CharacterModel) Insert(character *Character) error {
	// Insert a new character item into the database.
	query := `
		INSERT INTO characters (ID, FirstName, LastName, House, OriginStatus) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{character.ID, character.FirstName, character.LastName, character.House, character.OriginStatus}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&character.ID, &character.CreatedAt, &character.UpdatedAt)
}

func (c CharacterModel) Get(id int) (*Character, error) {
	// Retrieve a character item based on its ID.
	query := `
		SELECT ID, created_at, updated_at, FirstName, LastName, House, OriginStatus
		FROM characters
		WHERE id = $1
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
		SET firstName = $1, lastName = $2, House = $3
		WHERE ID = $4
		RETURNING updated_at
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