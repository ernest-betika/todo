package services

import (
	"context"
	"fmt"
	"strings"
	"time"
	"todo/internal/db"
	"todo/internal/entities"
	"todo/internal/forms"
	"todo/internal/repos"
)

type (
	TodoController interface {
		CompleteTodo(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error)
		CreateTodo(ctx context.Context, dB db.DB, form *forms.CreateTodoForm) (*entities.Todo, error)
		DeleteTodo(ctx context.Context, todoID int64) error
		TodoByID(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error)
		UpdateTodo(ctx context.Context, dB db.DB, todoID int64, form *forms.UpdateTodoForm) (*entities.Todo, error)
	}

	todoController struct {
		todoRepository repos.TodoRepository
	}
)

func NewTestTodoController(db db.DB) *todoController {
	return &todoController{
		todoRepository: repos.NewTodoRepository(db),
	}
}

func NewTodoController(
	todoRepository repos.TodoRepository,
) TodoController {
	return &todoController{
		todoRepository: todoRepository,
	}
}

func (s *todoController) TodoByID(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error) {

	todo, err := s.todoRepository.TodoByID(ctx, todoID)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoController) CreateTodo(ctx context.Context, dB db.DB, form *forms.CreateTodoForm) (*entities.Todo, error) {

	todo := &entities.Todo{
		Title: form.Title,
	}

	if strings.TrimSpace(form.Description) != "" {
		todo.Description = form.Description
	}

	err := s.todoRepository.Save(ctx, dB, todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoController) UpdateTodo(ctx context.Context, dB db.DB, todoID int64, form *forms.UpdateTodoForm) (*entities.Todo, error) {

	todo, err := s.todoRepository.TodoByID(ctx, todoID)
	if err != nil {
		return &entities.Todo{}, err
	}

	if form.Title != nil {
		todo.Title = *form.Title
	}

	if form.Description != nil {
		if strings.TrimSpace(*form.Description) != "" {
			todo.Description = *form.Description
		}
	}

	err = s.todoRepository.Save(ctx, dB, todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoController) CompleteTodo(ctx context.Context, dB db.DB, todoID int64) (*entities.Todo, error) {

	todo, err := s.todoRepository.TodoByID(ctx, todoID)
	if err != nil {
		return &entities.Todo{}, err
	}

	if todo.Completed {
		return &entities.Todo{}, fmt.Errorf("todo has been marked as complete")
	}

	timeNow := time.Now()

	todo.Completed = true
	todo.CompletedAt = &timeNow

	err = s.todoRepository.Save(ctx, dB, todo)
	if err != nil {
		return &entities.Todo{}, err
	}

	return todo, nil
}

func (s *todoController) DeleteTodo(ctx context.Context, todoID int64) error {

	todo, err := s.todoRepository.TodoByID(ctx, todoID)
	if err != nil {
		return err
	}

	if todo.Completed {
		return fmt.Errorf("cannot a todo that has been completed")
	}

	return s.todoRepository.DeleteTodo(ctx, todo.ID)
}
