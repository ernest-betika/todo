package todocomment

import (
	"fmt"
	"net/http"
	"strconv"
	"todo/internal/db"
	"todo/internal/forms"
	"todo/internal/services"

	"github.com/gin-gonic/gin"
)

func createTodoComment(
	dB db.DB,
	todoCommentService services.TodoCommentService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var form forms.CreateTodoCommentForm

		err := c.BindJSON(&form)
		if err != nil {
			return
		}

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return
		}

		ctx := c.Request.Context()

		//uncomment below to simulate the rollback
		// ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		// defer cancel()

		todoComment, err := todoCommentService.CreateTodoComment(ctx, dB, todoID, &form)
		if err != nil {
			fmt.Printf("todo comment service err %v", err)
			return
		}

		c.JSON(http.StatusOK, todoComment)
	}
}
