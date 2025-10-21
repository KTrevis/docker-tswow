package api

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type signupRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type signupResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func ensureUnique(db *sql.DB, usernameUpper, email string) error {
	var cnt int
	if err := db.QueryRow("SELECT COUNT(1) FROM account WHERE username = ?", usernameUpper).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("username already exists")
	}
	cnt = 0
	if err := db.QueryRow("SELECT COUNT(1) FROM account WHERE email = ?", email).Scan(&cnt); err != nil {
		return err
	}
	if cnt > 0 {
		return errors.New("email already exists")
	}
	return nil
}

// insertAccount inserts a new row into auth.account with generated salt/verifier.
func insertAccount(db *sql.DB, usernameUpper, email, password string) (int64, error) {
	// Placeholder salt/verifier. Replace with SRP-6 generation using username/password.
	salt, err := randomBytes(32)
	if err != nil {
		return 0, err
	}
	verifier, err := randomBytes(32)
	if err != nil {
		return 0, err
	}

	res, err := db.Exec(`
		INSERT INTO account (
			username, salt, verifier, email, reg_mail
		) VALUES (
			?, ?, ?, ?, ?
		)
	`, usernameUpper, salt, verifier, email, email)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// CreateAccount performs validation, uniqueness checks and DB insert. Useful for tests without HTTP.
func CreateAccount(db *sql.DB, username, email, password string) (int64, string, error) {
	usernameUpper := strings.ToUpper(strings.TrimSpace(username))
	if usernameUpper == "" {
		return 0, "", errors.New("invalid username")
	}
	if err := ensureUnique(db, usernameUpper, email); err != nil {
		return 0, "", err
	}
	id, err := insertAccount(db, usernameUpper, email, password)
	if err != nil {
		return 0, "", err
	}
	return id, usernameUpper, nil
}

func HandleSignup(c *gin.Context, db *sql.DB) {
	var body signupRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, usernameUpper, err := CreateAccount(db, body.Username, body.Email, body.Password)
	if err != nil {
		status := http.StatusInternalServerError
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, signupResponse{ID: id, Username: usernameUpper, Email: body.Email})
}
