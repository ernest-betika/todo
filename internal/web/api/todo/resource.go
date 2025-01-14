package todo

import (
	"net/http"
	"strconv"
	"todo/internal/db"
	"todo/internal/forms"
	"todo/internal/services"

	"github.com/gin-gonic/gin"
)

func createTodo(
	dB db.DB,
	todoController services.TodoController,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var form forms.CreateTodoForm

		err := c.BindJSON(&form)
		if err != nil {
			return
		}

		todo, err := todoController.CreateTodo(c.Request.Context(), dB, &form)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func completeTodo(
	dB db.DB,
	todoController services.TodoController,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return
		}

		todo, err := todoController.CompleteTodo(c.Request.Context(), dB, todoID)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func updateTodo(
	dB db.DB,
	todoController services.TodoController,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var form forms.UpdateTodoForm

		err := c.BindJSON(&form)
		if err != nil {
			return
		}

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return
		}

		todo, err := todoController.UpdateTodo(c.Request.Context(), dB, todoID, &form)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func todoByID(
	dB db.DB,
	todoController services.TodoController,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return
		}

		todo, err := todoController.TodoByID(c.Request.Context(), dB, todoID)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

func deleteTodo(
	todoController services.TodoController,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		todoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			return
		}

		err = todoController.DeleteTodo(c.Request.Context(), todoID)
		if err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	}
}
