package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Comment struct {
	Id          int    `json:"Id"`
	UsernameID  int64  `json:"UsernameID"`
	Comment     string `json:"Comment"`
	CharacterID int64  `json:"CharacterID"`
}

type CommentModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// CreateComment создает новый комментарий.
func (m *CommentModel) CreateComment(comment *Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO comments (UsernameID, Comment, CharacterID)
		VALUES ($1, $2, $3)
		RETURNING Id
	`

	err := m.DB.QueryRowContext(ctx, query, comment.UsernameID, comment.Comment, comment.CharacterID).Scan(&comment.Id)
	if err != nil {
		m.ErrorLog.Printf("Error inserting comment into database: %v", err)
		return err
	}
	

	return nil
}

// GetCommentByID возвращает комментарий по его ID.
func (m *CommentModel) GetCommentByID(commentID int) (*Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT Id, UsernameID, Comment, CharacterID
		FROM comments
		WHERE Id = $1
	`

	comment := &Comment{}

	err := m.DB.QueryRowContext(ctx, query, commentID).Scan(&comment.Id, &comment.UsernameID, &comment.Comment, &comment.CharacterID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Комментарий не найден
		}
		m.ErrorLog.Printf("Error querying comment from database: %v", err)
		return nil, err
	}

	return comment, nil
}

func (m *CommentModel) UpdateComment(comment *Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        UPDATE comments
        SET Comment = $1
        WHERE Id = $2
    `

	_, err := m.DB.ExecContext(ctx, query, comment.Comment, comment.Id)
	if err != nil {
		m.ErrorLog.Printf("Error updating comment in database: %v", err)
		return err
	}

	return nil
}

// DeleteCommentByID удаляет комментарий по его ID.
func (m *CommentModel) DeleteCommentByID(commentID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		DELETE FROM comments
		WHERE Id = $1
	`

	_, err := m.DB.ExecContext(ctx, query, commentID)
	if err != nil {
		m.ErrorLog.Printf("Error deleting comment from database: %v", err)
		return err
	}

	return nil
}

//фильтрация,сортировка,пагинация

// фильтр по айди юзера
func (m *CommentModel) GetCommentsByUserID(userID int64) ([]*Comment, error) {

	query := `
        SELECT Id, UsernameID, Comment, CharacterID
        FROM comments
        WHERE UsernameID = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		err := rows.Scan(&comment.Id, &comment.UsernameID, &comment.Comment, &comment.CharacterID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

//сортировка по айди персонажа
func (m *CommentModel) GetCommentsByCharacter() ([]*Comment, error) {

    query := `
        SELECT Id, UsernameID, Comment, CharacterID
        FROM comments
        ORDER BY CharacterID 
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    rows, err := m.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []*Comment
    for rows.Next() {
        comment := &Comment{}
        err := rows.Scan(&comment.Id, &comment.UsernameID, &comment.Comment, &comment.CharacterID)
        if err != nil {
            return nil, err
        }
        comments = append(comments, comment)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return comments, nil
}

//пагинация
func (m *CommentModel) GetCommentsPagination(limit, offset int) ([]*Comment, error) {
    
    query := `
        SELECT Id, UsernameID, Comment, CharacterID
        FROM comments
        LIMIT $1 OFFSET $2
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    rows, err := m.DB.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []*Comment
    for rows.Next() {
        comment := &Comment{}
        err := rows.Scan(&comment.Id, &comment.UsernameID, &comment.Comment, &comment.CharacterID)
        if err != nil {
            return nil, err
        }
        comments = append(comments, comment)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return comments, nil
}

func (m *CommentModel) GetCommentsByCharacterID(characterID int) ([]*Comment, error) {
	query := `
		SELECT Id, UsernameID, Comment, CharacterID
		FROM comments
		WHERE CharacterID = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, characterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		err := rows.Scan(&comment.Id, &comment.UsernameID, &comment.Comment, &comment.CharacterID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// выводим список комментов от определенного юзера
func (m *CommentModel) GetCommentsByUser(userID int) ([]*Comment, error) {
	query := `
		SELECT Id, UsernameID, Comment, CharacterID
		FROM comments
		WHERE UsernameID = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		comment := &Comment{}
		err := rows.Scan(&comment.Id, &comment.UsernameID, &comment.Comment, &comment.CharacterID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

