package router

import (
	"net/http"
	"os"
	"todo/internal/db"
	"todo/internal/repos"
	"todo/internal/services"
	"todo/internal/web/api/todo"
	todocomment "todo/internal/web/api/todo_comment"

	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	*gin.Engine
}

func BuildRouter(
	dB db.DB,
) *AppRouter {

	if os.Getenv("ENVIRONMENT") == "development" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	appRouter := router.Group("/v1")

	todoRepository := repos.NewTodoRepository(dB)
	todoCommentRepository := repos.NewTodoCommentRepository(dB)
	transactionProviderRepository := repos.NewTransactionProviderRepository(dB)

	todoController := services.NewTodoController(todoRepository)
	todoCommentService := services.NewTodoCommentService(
		todoRepository,
		todoCommentRepository,
		transactionProviderRepository,
	)

	todo.AddOpenEndpoints(appRouter, dB, todoController)
	todocomment.AddOpenEndpoints(appRouter, dB, todoCommentService)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Endpoint not found"})
	})

	return &AppRouter{
		router,
	}
}
