package storage

import (
	"context"
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
	"github.com/GermanBogatov/TodoApp/app/pkg/postgresql"
)

const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	GetUser(ctx context.Context, username, password string) (model.User, error)
}

type TodoList interface {
	Create(ctx context.Context, userId int, list model.TodoList) (int, error)
	GetAll(ctx context.Context, userId int) ([]model.TodoList, error)
	GetById(ctx context.Context, userId, listId int) (model.TodoList, error)
	Delete(ctx context.Context, userId, listId int) error
	Update(ctx context.Context, userId, listId int, input model.UpdateListInput) error
}

type TodoItem interface {
	Create(ctx context.Context, listId int, item model.TodoItem) (int, error)
	GetAll(ctx context.Context, userId, listId int) ([]model.TodoItem, error)
	GetById(ctx context.Context, userId, itemId int) (model.TodoItem, error)
	Delete(ctx context.Context, userId, itemId int) error
	Update(ctx context.Context, userId, itemId int, input model.UpdateItemInput) error
}

type Storage struct {
	Authorization
	TodoList
	TodoItem
}

//constructor
func NewRepository(client postgresql.ClientPostgres, logger logging.Logger) *Storage {
	return &Storage{
		Authorization: NewAuthorization(client, logger),
		TodoList:      NewLists(client, logger),
		TodoItem:      NewItems(client, logger),
	}
}
