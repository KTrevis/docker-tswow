package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(apiGroup *gin.RouterGroup, db *sql.DB) {
	handleSignup := func(c *gin.Context) { HandleSignup(c, db) }

	apiGroup.POST("/signup", handleSignup)
}
