package api

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(apiGroup *gin.RouterGroup, db *sql.DB) {
	handleSignup := func(c *gin.Context) { HandleSignup(c, db) }

	log.Println("Creating account...")
	id, username, err := CreateAccount(db, "test", "test@test.com", "password")
	log.Println("id", id, "username", username, "err", err)
	apiGroup.POST("/signup", handleSignup)
}
