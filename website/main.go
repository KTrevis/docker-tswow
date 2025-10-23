package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"website/api"
	"website/db"
)

func main() {
	database, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}
	if err := database.Ping(); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world!")
	})

	r.GET("/register", func(c *gin.Context) {
		c.File("static/register.html")
	})

	apiGroup := r.Group("/api")
	api.RegisterRoutes(apiGroup, database)

	log.Println("Listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
