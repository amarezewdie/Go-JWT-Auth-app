package models

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents the structure of a user in the system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// get user with role
func GetUserWithRole(db *sql.DB, id int) (*User, error) {
	query := `SELECT id, name, email, role, created_at, updated_at FROM users WHERE id = ?`
	row := db.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
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
		return fmt.Errorf("hash password failed: %w", err)
	}

	query := `
		INSERT INTO users (name, email, password, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query, user.Name, user.Email, user.Password, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert user failed: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last inserted id failed: %w", err)
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
		SELECT id, name, email, password, role, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	row := db.QueryRow(query, email)

	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user from DB using their ID
func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := `
		SELECT id, name, email, role, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	row := db.QueryRow(query, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("get user by ID failed: %w", err)
	}
	return user, nil
}

// GetAllUsers returns all users from the database
func GetAllUsers(db *sql.DB) ([]User, error) {
	query := `
		SELECT id, name, email, role, created_at, updated_at
		FROM users
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("get all users query failed: %w", err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
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

// InitAdminUser inserts the default admin into the database if not already present
func InitAdminUser(db *sql.DB) error {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminName := os.Getenv("ADMIN_NAME")

	if adminEmail == "" || adminPassword == "" || adminName == "" {
		return fmt.Errorf("admin credentials are missing in environment variables")
	}

	// Check if admin already exists
	existingUser, err := GetUserByEmail(db, adminEmail)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check for existing admin: %w", err)
	}
	if existingUser != nil {
		fmt.Println("✅ Admin already exists:", existingUser.Email)
		return nil // already created
	}

	// Create new admin user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing admin password failed: %w", err)
	}

	now := time.Now()

	query := `
		INSERT INTO users (name, email, password, role, created_at, updated_at)
		VALUES (?, ?, ?, 'admin', ?, ?)
	`

	_, err = db.Exec(query, adminName, adminEmail, hashedPassword, now, now)
	if err != nil {
		return fmt.Errorf("failed to insert admin user: %w", err)
	}

	fmt.Println("🎉 Admin user created:", adminEmail)
	return nil
}
