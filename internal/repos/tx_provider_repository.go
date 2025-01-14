package repos

import (
	"todo/internal/db"
)

type TransactionProvider interface {
	InTransaction(txFunc func(adapters TransactionRepository) error) error
}

type transactionProvider struct {
	operations db.SQLOperations
}

type TransactionRepository struct {
	TodoRepository        *todoRepository
	TodoCommentRepository *todoCommentRepository
}

func NewTransactionProviderRepository(operations db.SQLOperations) TransactionProvider {
	return &transactionProvider{operations: operations}
}

func (r *transactionProvider) InTransaction(txProviders func(transactionRepository TransactionRepository) error) error {
	return db.WithTransaction(db.GetDB(), func(operations db.SQLOperations) error {
		transactionRepository := TransactionRepository{
			TodoRepository:        NewTodoRepository(operations),
			TodoCommentRepository: NewTodoCommentRepository(operations),
		}

		return txProviders(transactionRepository)
	})
}
