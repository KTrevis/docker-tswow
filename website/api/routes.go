package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(apiGroup *gin.RouterGroup, db *sql.DB) {
	handleSignup := func(c *gin.Context) { HandleSignup(c, db) }

	CreateAccount(db, "test", "test@test.com", "password")
	apiGroup.POST("/signup", handleSignup)
}
