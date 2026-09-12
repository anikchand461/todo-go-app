package routes

import (
	"github.com/gin-gonic/gin"

	"todo-go/backend/auth"
	"todo-go/backend/handlers"
)

func SetupRoutes(
	router *gin.Engine,
	todoHandler *handlers.TodoHandler,
	authHandler *handlers.AuthHandler,
) {
	api := router.Group("/api")

	// Public authentication routes.
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	// Protected Todo routes.
	todos := api.Group("/todos")
	todos.Use(auth.Middleware())

	todos.GET("", todoHandler.GetTodos)
	todos.GET("/:id", todoHandler.GetTodo)

	todos.POST("", todoHandler.CreateTodo)

	todos.PUT("/:id", todoHandler.UpdateTodo)

	todos.DELETE("/:id", todoHandler.DeleteTodo)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Todo API is running",
		})
	})
}
