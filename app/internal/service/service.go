package service

import (
	"context"
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/GermanBogatov/TodoApp/app/internal/storage"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(ctx context.Context, token string) (int, error)
}

type TodoList interface {
	Create(ctx context.Context, userId int, list model.TodoList) (int, error)
	GetAll(ctx context.Context, userId int) ([]model.TodoList, error)
	GetById(ctx context.Context, userId, listId int) (model.TodoList, error)
	Delete(ctx context.Context, userId, listId int) error
	Update(ctx context.Context, userId, listId int, input model.UpdateListInput) error
}

type TodoItem interface {
	Create(ctx context.Context, userId, listId int, item model.TodoItem) (int, error)
	GetAll(ctx context.Context, userId, listId int) ([]model.TodoItem, error)
	GetById(ctx context.Context, userId, itemId int) (model.TodoItem, error)
	Delete(ctx context.Context, userId, itemId int) error
	Update(ctx context.Context, userId, itemId int, input model.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(storage *storage.Storage, logger logging.Logger) (*Service, error) {
	return &Service{
		Authorization: NewAuthService(storage.Authorization, logger),
		TodoList:      NewTodoListService(storage.TodoList, logger),
		TodoItem:      NewTodoItemService(storage.TodoItem, storage.TodoList, logger),
	}, nil
}
