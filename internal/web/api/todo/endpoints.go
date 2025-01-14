package todo

import (
	"todo/internal/db"
	"todo/internal/services"

	"github.com/gin-gonic/gin"
)

func AddOpenEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	todoController services.TodoController,
) {
	r.POST("/todo", createTodo(dB, todoController))
	r.GET("/todo/:id", todoByID(dB, todoController))
	r.PUT("/todo/:id", updateTodo(dB, todoController))
	r.POST("/todo/:id", completeTodo(dB, todoController))
	r.DELETE("/todo/:id", deleteTodo(todoController))
}
