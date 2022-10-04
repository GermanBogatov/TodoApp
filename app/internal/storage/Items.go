package storage

import (
	"context"
	"fmt"
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
	"github.com/GermanBogatov/TodoApp/app/pkg/postgresql"
	"github.com/jackc/pgconn"
	"strings"
)

type repositoryItems struct {
	client postgresql.ClientPostgres
	logger logging.Logger
}

func NewItems(client postgresql.ClientPostgres, logger logging.Logger) *repositoryItems {
	return &repositoryItems{
		client: client,
		logger: logger,
	}
}

func (r *repositoryItems) Create(ctx context.Context, listId int, item model.TodoItemDTO) (int, error) {
	tx, err := r.client.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer tx.Commit(ctx)
	//  todoItemsTable  = "todo_items"
	//	listsItemsTable = "lists_items"

	qitems := `
	INSERT INTO todo_items
		(title,description)
	VALUES
		($1,$2)
	RETURNING id
	`
	if err = tx.QueryRow(ctx, qitems, item.Title, item.Description).Scan(&item.Id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return 0, newErr
		}
		return 0, err
	}

	qlists := `
	INSERT INTO lists_items
		(item_id, list_id)
	VALUES
		($1,$2)
	`
	_, err = tx.Exec(ctx, qlists, item.Id, listId)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	return item.Id, nil

}

func (r *repositoryItems) GetAll(ctx context.Context, userId, listId int) ([]model.TodoItemDTO, error) {
	r.logger.Println("GET ALL ITEMS")
	q := `
	SELECT 
		ti.id, ti.title, ti.description, ti.done
	FROM
		todo_items ti
	INNER JOIN
		lists_items li on li.item_id = ti.id
	INNER JOIN
		users_lists ul on ul.list_id = li.list_id
	WHERE
		li.list_id = $1
	AND
		ul.user_id = $2
		`
	rows, err := r.client.Query(ctx, q, listId, userId)

	if err != nil {
		return nil, err
	}
	items := make([]model.TodoItemDTO, 0)

	for rows.Next() {
		var item model.TodoItemDTO

		err = rows.Scan(&item.Id, &item.Title, &item.Description, &item.Done)
		if err != nil {
			return nil, err
		}

		items = append(items, item)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return items, nil

}

func (r *repositoryItems) GetById(ctx context.Context, userId, itemId int) (model.TodoItemDTO, error) {
	var item model.TodoItemDTO

	q := `
	SELECT
		ti.id, ti.title, ti.description, ti.done
	FROM
		todo_items ti
	INNER JOIN
		lists_items li on li.item_id = ti.id
	INNER JOIN 
		users_lists ul on ul.list_id = li.list_id
	WHERE 
		ti.id = $1 
	AND 
		ul.user_id = $2
`
	err := r.client.QueryRow(ctx, q, itemId, userId).Scan(&item.Id, &item.Title, &item.Description, &item.Done)
	if err != nil {
		return model.TodoItemDTO{}, err
	}

	return item, nil
}

func (r *repositoryItems) Delete(ctx context.Context, userId, itemId int) error {
	q := `
	DELETE FROM
		todo_items ti
	USING 
		lists_items li, users_lists ul
	WHERE 
		ti.id = li.item_id 
	AND 
		li.list_id = ul.list_id 
	AND 
		ul.user_id = $1 
	AND 
		ti.id = $2
`
	_, err := r.client.Exec(ctx, q, userId, itemId)
	return err
}

func (r *repositoryItems) Update(ctx context.Context, userId, itemId int, input model.UpdateItemInput) error {

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul
									WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)

	_, err := r.client.Exec(ctx, query, args...)
	return err

}
