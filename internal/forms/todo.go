package forms

type CreateTodoForm struct {
	Description string `json:"description" binding:"required"`
	Title       string `json:"title" binding:"required"`
}

type UpdateTodoForm struct {
	Description *string `json:"description"`
	Title       *string `json:"title"`
}