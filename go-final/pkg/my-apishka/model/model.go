package model

import (
	"database/sql"
	"log"
	"os"
)


type Models struct {
	Characters CharacterModel
	Users UserModel 
	Tokens TokenModel
	Permissions PermissionModel
	Comments CommentModel
}


func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		Characters: CharacterModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Users: UserModel{
			DB: db,
		},
		Tokens: TokenModel{
			DB: db,
		},
		Permissions: PermissionModel{
			DB: db,
		},
		Comments: CommentModel{
			DB: db,
		},
	}
}