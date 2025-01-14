package todocomment

import (
	"todo/internal/db"
	"todo/internal/services"

	"github.com/gin-gonic/gin"
)

func AddOpenEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	todoCommentService services.TodoCommentService,
) {
	r.POST("/todo/:id/todo_comment", createTodoComment(dB, todoCommentService))
}
