package repos

import (
	"context"
	"fmt"
	"time"
	"todo/internal/db"
	"todo/internal/entities"
)

const (
	insertTodoCommentSQL = "INSERT INTO todos_comments (comment, todo_id, created_at) VALUES ($1, $2, $3) RETURNING id"
)

type (
	TodoCommentRepository interface {
		Save(ctx context.Context, operations db.SQLOperations, todoComment *entities.TodoComment) error
		NSave(ctx context.Context, todoComment *entities.TodoComment) error
	}

	todoCommentRepository struct {
		operations db.SQLOperations
	}
)

func NewTodoCommentRepository(
	operations db.SQLOperations,
) TodoCommentRepository {
	return &todoCommentRepository{
		operations: operations,
	}
}

func (r *todoCommentRepository) NSave(ctx context.Context, todoComment *entities.TodoComment) error {

	if todoComment.IsNew() {

		todoComment.CreatedAt = time.Now()

		err := r.operations.QueryRowContext(
			ctx,
			insertTodoCommentSQL,
			todoComment.Comment,
			todoComment.TodoID,
			todoComment.CreatedAt,
		).Scan(&todoComment.ID)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("cannot update todo comment")
}

func (r *todoCommentRepository) Save(ctx context.Context, operations db.SQLOperations, todoComment *entities.TodoComment) error {

	if todoComment.IsNew() {

		todoComment.CreatedAt = time.Now()

		err := operations.QueryRowContext(
			ctx,
			insertTodoCommentSQL,
			todoComment.Comment,
			todoComment.TodoID,
			todoComment.CreatedAt,
		).Scan(&todoComment.ID)
		if err != nil {
			fmt.Printf("todo repo err %v", err)
			return err
		}

		return nil
	}

	return fmt.Errorf("cannot update todo comment")
}
