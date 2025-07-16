package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents the structure of a user in the system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // hide password in JSON response
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

///////////////////////////////////////////////////////////////////////////////
// 🔐 PASSWORD UTILITIES
///////////////////////////////////////////////////////////////////////////////

// HashPassword hashes the user's password using bcrypt before saving to DB
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// ComparePassword compares raw input password with hashed one stored in DB
func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

///////////////////////////////////////////////////////////////////////////////
// 🟢 CREATE USER
///////////////////////////////////////////////////////////////////////////////

// CreateUser inserts a new user into the database
func CreateUser(db *sql.DB, user *User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := user.HashPassword(); err != nil {
		return err
	}

	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	// Capture the auto-incremented ID from the DB
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// 🔍 READ USERS
///////////////////////////////////////////////////////////////////////////////

// GetUserByEmail retrieves a user from DB using their email
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	row := db.QueryRow(query, email)

	user := &User{}
	err := row.Scan(
		&user.ID, &user.Name, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user from DB using their ID
func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	row := db.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAllUsers returns all users from the database
func GetAllUsers(db *sql.DB) ([]User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

///////////////////////////////////////////////////////////////////////////////
// ✏️ UPDATE USER
///////////////////////////////////////////////////////////////////////////////

// UpdateUser modifies user details in the DB (excluding password)
func UpdateUser(db *sql.DB, id int, user *User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET name = ?, email = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := db.Exec(query, user.Name, user.Email, user.UpdatedAt, id)
	return err
}

///////////////////////////////////////////////////////////////////////////////
// ❌ DELETE USER
///////////////////////////////////////////////////////////////////////////////

// DeleteUser removes a user from the database by ID
func DeleteUser(db *sql.DB, id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}
