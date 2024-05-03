package model

import (
	"database/sql"
	"errors"
	"time"

	"go-final/pkg/my-apishka/validator"
	"golang.org/x/crypto/bcrypt"

	"context"

	"gorm.io/gorm"

	"crypto/sha256"
)


type User struct {
	ID         int64     `json:"ID"`
	CreatedAt  time.Time `json:"CreatedAt"`
	Username   string    `json:"Username"`
	Email      string    `json:"Email"`
	Password   password  `json:"-"`
	Activated  bool      `json:"Activated"`
	Version    int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var(
	ErrEditConflict = errors.New("edit conflict")
)

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}


func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}




type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (Username, Email, Password, Activated)
		VALUES ($1, $2, $3, $4)
		RETURNING ID, CreatedAt, Version
		`

	args := []interface{}{user.Username, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// If the table already contains a record with this email address, then when we try to
	// perform the insert there will be a violation of the UNIQUE "users_email_key" constraint
	// that we set up in the previous chapter. We check for this error specifically, and return
	// ErrDuplicateEmail error instead.
	pqErr := `pq: duplicate key value violates unique constraint "users_email_key"`
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == pqErr:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT ID, CreatedAt, Username, Email, Password, Activated, Version
		FROM users
		WHERE Email = $1
		`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, gorm.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET Username = $1, Email = $2, Password = $3, Activated = $4, Version = Version + 1
		WHERE ID = $5 AND Version = $6
		RETURNING Version
		`

	args := []interface{}{
		user.Username,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	// validate user.username
	v.Check(user.Username != "", "name", "must be provided")
	v.Check(len(user.Username) <= 500, "name", "must not be more than 500 bytes long")

	// Validate email
	ValidateEmail(v, user.Email)

	// If the plaintext password is not nil, call the standalone ValidatePasswordPlaintext helper.
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	// If the password has is ever nil, this will be due to a logic error in our codebase
	// (probably because we forgot to set a password for the user). It's a useful sanity check to
	// include here, but it's not a problem with the data provided by the client. So, rather
	// than adding an error to the validation map we raise a panic instead.
	if user.Password.hash == nil {
		// TODO: fix this panic
		panic("missing password hash for user")
	}
}


func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	// Calculate the SHA-256 hash for the plaintext token provided by the client.
	// Note, that this will return a byte *array* with length 32, not a slice.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
		SELECT 
			users.ID, users.CreatedAt, users.Username, users.Email, 
			users.Password, users.Activated, users.Version
		FROM       users
        INNER JOIN tokens
			ON users.id = tokens.user_id
        WHERE tokens.hash = $1  --<-- Note: this is potentially vulnerable to a timing attack, 
            -- but if successful the attacker would only be able to retrieve a *hashed* token 
            -- which would still require a brute-force attack to find the 26 character string
            -- that has the same SHA-256 hash that was found from our database. 
			AND tokens.scope = $2
			AND tokens.expiry > $3
		`

	// Create a slice containing the query args. Note, that we use the [:] operator to get a slice
	// containing the token hash, since the pq driver does not support passing in an array.
	// Also, we pass the current time as the value to check against the token expiry.
	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query, scanning the return values into a User struct. If no matching record
	// is found we return an ErrRecordNotFound error.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, gorm.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Return the matching user.
	return &user, nil
}