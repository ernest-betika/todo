package services

import (
	"context"
	"fmt"
	"time"
	"todo/internal/db"
	"todo/internal/entities"
	"todo/internal/forms"
	"todo/internal/repos"
)

type (
	TodoCommentService interface {
		CreateTodoComment(ctx context.Context, dB db.DB, todoID int64, form *forms.CreateTodoCommentForm) (*entities.TodoComment, error)
	}

	todoCommentService struct {
		todoRepository        repos.TodoRepository
		todoCommentRepository repos.TodoCommentRepository
		txHelper1             repos.TransactionProvider
		txHelper2             repos.TransactionProvider
	}
)

func NewTodoCommentService(
	todoRepository repos.TodoRepository,
	todoCommentRepository repos.TodoCommentRepository,
	txHelper1 repos.TransactionProvider,
	txHelper2 repos.TransactionProvider,
) TodoCommentService {
	return &todoCommentService{
		todoRepository:        todoRepository,
		todoCommentRepository: todoCommentRepository,
		txHelper1:             txHelper1,
		txHelper2:             txHelper2,
	}
}

func (s *todoCommentService) CreateTodoComment(
	ctx context.Context,
	dB db.DB,
	todoID int64,
	form *forms.CreateTodoCommentForm,
) (*entities.TodoComment, error) {

	todo, err := s.todoRepository.TodoByID(ctx, todoID)
	if err != nil {
		return &entities.TodoComment{}, err
	}

	todoComment := &entities.TodoComment{
		Comment: form.Comment,
		TodoID:  todo.ID,
	}

	/*
		this is the first option to simulate rollback usage is by passing db connection as a parameter

	*/
	/*
		err = dB.InTransaction(ctx, func(ctx context.Context, operations db.SQLOperations) error {

			err = s.todoCommentRepository.Save(ctx, operations, todoComment)
			if err != nil {
				fmt.Printf("todo comment err %v", err.Error())
				return err
			}

			time.Sleep(5 * time.Second)

			todo.Completed = true

			err := s.todoRepository.Save(ctx, operations, todo)
			if err != nil {
				fmt.Printf("todo err %v", err.Error())
				return err
			}

			return nil
		})
	*/

	err = s.txHelper1.InTransaction(ctx, func(transactionRepository repos.TransactionRepository) error {

		todo.Completed = true

		err := transactionRepository.TodoRepository.NSave(ctx, todo)
		if err != nil {
			fmt.Printf("todo err %v", err)
			return err
		}

		err = transactionRepository.TodoCommentRepository.NSave(ctx, todoComment)
		if err != nil {
			fmt.Printf("todo comment err %v", err)
			return err
		}

		err = s.txHelper2.InTransaction(ctx, func(innerTransactionRepository repos.TransactionRepository) error {

			todo1, err := innerTransactionRepository.TodoRepository.TodoByID(ctx, 19)
			if err != nil {
				return err
			}

			todo1.Completed = true

			err = innerTransactionRepository.TodoRepository.NSave(ctx, todo1)
			if err != nil {
				fmt.Printf("inner todo err %v", err)
				return err
			}

			//simulate sleep. If context.Timeout is set error will be printed of context exceeded and the above query will be rolled back.
			// check on todo repository completed column
			time.Sleep(5 * time.Second)

			todoComment.ID = 0
			todoComment.TodoID = todo1.ID

			err = innerTransactionRepository.TodoCommentRepository.NSave(ctx, todoComment)
			if err != nil {
				fmt.Printf("inner todo comment err %v", err)
				return err
			}

			return nil
		})

		return err
	})
	if err != nil {
		return nil, err
	}

	return todoComment, nil
}
