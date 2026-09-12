package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"todo-go/backend/config"
	"todo-go/backend/database"
	"todo-go/backend/handlers"
	"todo-go/backend/routes"
)

func main() {
	cfg := config.Load()

	db := database.Connect(cfg.DatabaseURL)
	defer db.Close()

	router := gin.Default()

	todoHandler := handlers.NewTodoHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	routes.SetupRoutes(
		router,
		todoHandler,
		authHandler,
	)

	log.Printf(
		"Todo API running on http://localhost:%s",
		cfg.Port,
	)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
