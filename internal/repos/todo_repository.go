package repos

import (
	"context"
	"todo/internal/db"
	"todo/internal/entities"
)

const (
	deleteTodoSQL  = "DELETE FROM todos WHERE id = $1"
	getTodoByIDSQL = selectTodoSQL + " WHERE id = $1"
	insertTodoSQL  = "INSERT INTO todos (title, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id"
	selectTodoSQL  = "SELECT id, title, description, completed, completed_at, created_at, updated_at FROM todos"
	updateTodoSQL  = "UPDATE todos SET title = $1, description = $2, completed = $3, completed_at = $4, updated_at = $5 WHERE id = $6"
)

type (
	TodoRepository interface {
		DeleteTodo(ctx context.Context, todoID int64) error
		//below to methods are similar only difference is that one passes the db connection as a param
		Save(ctx context.Context, operations db.SQLOperations, todo *entities.Todo) error
		NSave(ctx context.Context, todo *entities.Todo) error
		
		TodoByID(ctx context.Context, todoID int64) (*entities.Todo, error)
		Todos(ctx context.Context) ([]*entities.Todo, error)
	}

	todoRepository struct {
		operations db.SQLOperations
	}
)

func NewTodoRepository(operations db.SQLOperations) TodoRepository {
	return &todoRepository{operations: operations}
}

func (r *todoRepository) NSave(
	ctx context.Context,
	todo *entities.Todo,
) error {

	todo.Touch()

	if todo.IsNew() {

		err := r.operations.QueryRowContext(
			ctx,
			insertTodoSQL,
			todo.Title,
			todo.Description,
			todo.CreatedAt,
			todo.UpdatedAt,
		).Scan(&todo.ID)
		if err != nil {
			return err
		}

		return nil
	}

	_, err := r.operations.ExecContext(
		ctx,
		updateTodoSQL,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.CompletedAt,
		todo.UpdatedAt,
		todo.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *todoRepository) Save(
	ctx context.Context,
	operations db.SQLOperations,
	todo *entities.Todo,
) error {

	todo.Touch()

	if todo.IsNew() {

		err := operations.QueryRowContext(
			ctx,
			insertTodoSQL,
			todo.Title,
			todo.Description,
			todo.CreatedAt,
			todo.UpdatedAt,
		).Scan(&todo.ID)
		if err != nil {
			return err
		}

		return nil
	}

	_, err := operations.ExecContext(
		ctx,
		updateTodoSQL,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.CompletedAt,
		todo.UpdatedAt,
		todo.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *todoRepository) TodoByID(
	ctx context.Context,
	todoID int64,
) (*entities.Todo, error) {

	row := r.operations.QueryRowContext(
		ctx,
		getTodoByIDSQL,
		todoID,
	)

	return r.scanRow(row)
}

func (r *todoRepository) Todos(
	ctx context.Context,
) ([]*entities.Todo, error) {

	query := selectTodoSQL
	args := []any{}

	rows, err := r.operations.QueryContext(ctx, query, args...)
	if err != nil {
		return []*entities.Todo{}, err
	}

	defer rows.Close()

	todos := make([]*entities.Todo, 0)

	for rows.Next() {
		todo, err := r.scanRow(rows)
		if err != nil {
			return []*entities.Todo{}, err
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return []*entities.Todo{}, err
	}

	return todos, nil
}

func (r *todoRepository) DeleteTodo(
	ctx context.Context,
	todoID int64,
) error {

	_, err := r.operations.ExecContext(
		ctx,
		deleteTodoSQL,
		todoID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *todoRepository) scanRow(
	rowScanner db.RowScanner,
) (*entities.Todo, error) {

	var todo entities.Todo

	err := rowScanner.Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CompletedAt,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		return &entities.Todo{}, err
	}

	return &todo, nil
}
