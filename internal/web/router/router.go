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
	dB1 db.DB,
	dB2 db.DB,
) *AppRouter {

	if os.Getenv("ENVIRONMENT") == "development" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	appRouter := router.Group("/v1")

	todoRepository := repos.NewTodoRepository(dB1)
	todoCommentRepository := repos.NewTodoCommentRepository(dB1)
	transactionProvider1Repository := repos.NewTransactionProviderRepository(dB1)
	transactionProvider2Repository := repos.NewTransactionProviderRepository(dB2)

	todoController := services.NewTodoController(todoRepository)
	todoCommentService := services.NewTodoCommentService(
		todoRepository,
		todoCommentRepository,
		transactionProvider1Repository,
		transactionProvider2Repository,
	)

	todo.AddOpenEndpoints(appRouter, dB1, todoController)
	todocomment.AddOpenEndpoints(appRouter, dB1, todoCommentService)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Endpoint not found"})
	})

	return &AppRouter{
		router,
	}
}
