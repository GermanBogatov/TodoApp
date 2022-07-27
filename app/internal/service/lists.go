package service

import (
	"context"
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/GermanBogatov/TodoApp/app/internal/storage"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
)

type listService struct {
	storageLists storage.TodoList
	logger       logging.Logger
}

func NewTodoListService(storageLists storage.TodoList, logger logging.Logger) *listService {
	return &listService{
		storageLists: storageLists,
		logger:       logger,
	}
}

func (s *listService) Create(ctx context.Context, userId int, list model.TodoList) (int, error) {
	return s.storageLists.Create(ctx, userId, list)
}

func (s *listService) GetAll(ctx context.Context, userId int) ([]model.TodoList, error) {
	return s.storageLists.GetAll(ctx, userId)
}

func (s *listService) GetById(ctx context.Context, userId, listId int) (model.TodoList, error) {
	return s.storageLists.GetById(ctx, userId, listId)
}

func (s *listService) Delete(ctx context.Context, userId, listId int) error {
	return s.storageLists.Delete(ctx, userId, listId)
}

func (s *listService) Update(ctx context.Context, userId, listId int, input model.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.storageLists.Update(ctx, userId, listId, input)
	return nil
}
