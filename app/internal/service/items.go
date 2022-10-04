package service

import (
	"context"
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/GermanBogatov/TodoApp/app/internal/storage"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
)

type itemService struct {
	storageItems storage.TodoItem
	storageLists storage.TodoList
	logger       logging.Logger
}

func NewTodoItemService(storageItems storage.TodoItem, storageLists storage.TodoList, logger logging.Logger) *itemService {
	return &itemService{
		storageItems: storageItems,
		storageLists: storageLists,
		logger:       logger,
	}
}

func (s *itemService) Create(ctx context.Context, userId, listId int, item model.TodoItemDTO) (int, error) {
	_, err := s.storageLists.GetById(ctx, userId, listId)
	if err != nil {
		//list does not exists or does not belongs to user
		return 0, err
	}

	return s.storageItems.Create(ctx, listId, item)
}

func (s *itemService) GetAll(ctx context.Context, userId, listId int) ([]model.TodoItemDTO, error) {
	return s.storageItems.GetAll(ctx, userId, listId)
}

func (s *itemService) GetById(ctx context.Context, userId, itemId int) (model.TodoItemDTO, error) {
	return s.storageItems.GetById(ctx, userId, itemId)
}

func (s *itemService) Delete(ctx context.Context, userId, itemId int) error {
	return s.storageItems.Delete(ctx, userId, itemId)
}

func (s *itemService) Update(ctx context.Context, userId, itemId int, input model.UpdateItemInput) error {
	return s.storageItems.Update(ctx, userId, itemId, input)
}
