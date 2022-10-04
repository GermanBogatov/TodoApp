package storage

import (
	"context"
	"fmt"
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
	"github.com/GermanBogatov/TodoApp/app/pkg/postgresql"
	"github.com/jackc/pgconn"
	"github.com/sirupsen/logrus"
	"strings"
)

type repositoryLists struct {
	client postgresql.ClientPostgres
	logger logging.Logger
}

func NewLists(client postgresql.ClientPostgres, logger logging.Logger) *repositoryLists {
	return &repositoryLists{
		client: client,
		logger: logger,
	}
}

func (r *repositoryLists) Create(ctx context.Context, userId int, list model.TodoListDTO) (int, error) {

	qLists := `
    INSERT INTO todo_lists 
    	(title, description) 
    VALUES 
		($1,$2) 
    RETURNING id
		`

	tx, err := r.client.Begin(ctx)
	if err != nil {
		return 0, nil
	}

	defer tx.Commit(ctx)

	if err = tx.QueryRow(ctx, qLists, list.Title, list.Description).Scan(&list.Id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return 0, newErr
		}
		return 0, err
	}
	qusers := `
	INSERT INTO users_lists
		(user_id,list_id)
	VALUES
		($1,$2)
	RETURNING id
	`

	_, err = tx.Exec(ctx, qusers, userId, list.Id)
	if err != nil {
		return 0, err
	}
	//fmt.Println("INFORMATION:: ", tx.Conn().ConnInfo())
	return list.Id, nil

}

func (r *repositoryLists) GetAll(ctx context.Context, userId int) ([]model.TodoListDTO, error) {

	q := `
		SELECT
			tl.id, tl.title, tl.description
		FROM
			todo_lists tl
		INNER JOIN
			users_lists ul on tl.id =ul.list_id
		WHERE ul.user_id = $1
	`
	rows, err := r.client.Query(ctx, q, userId)
	if err != nil {
		return nil, err
	}

	lists := make([]model.TodoListDTO, 0)

	for rows.Next() {
		var list model.TodoListDTO

		err = rows.Scan(&list.Id, &list.Title, &list.Description)
		if err != nil {
			return nil, err
		}

		lists = append(lists, list)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return lists, nil

}

func (r *repositoryLists) GetById(ctx context.Context, userId, listId int) (model.TodoListDTO, error) {

	var list model.TodoListDTO
	q := `
			SELECT
				tl.id, tl.title, tl.description
			FROM
				todo_lists tl
			INNER JOIN
				users_lists ul on tl.id = ul.list_id
			WHERE
				ul.user_id = $1
			AND
				ul.list_id = $2
		`
	err := r.client.QueryRow(ctx, q, userId, listId).Scan(&list.Id, &list.Title, &list.Description)
	if err != nil {
		return model.TodoListDTO{}, err
	}
	return list, nil

}

func (r *repositoryLists) Delete(ctx context.Context, userId, listId int) error {
	q := `
		DELETE FROM
			todo_lists tl
		USING
			users_lists ul
		WHERE
			tl.id = ul.list_id
		AND
			ul.user_id=$1
		AND
			ul.list_id=$2
	`
	_, err := r.client.Exec(ctx, q, userId, listId)
	return err

}

func (r *repositoryLists) Update(ctx context.Context, userId, listId int, input model.UpdateListInputDTO) error {

	/*q := `
		UPDATE
			todo_lists tl
		SET
			$1
		FROM
			users_lists ul
		WHERE
			tl.id = ul.list_id
		AND
			ul.list_id=$2
		AND
			ul.user_id=$3
	`*/

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

	// title=$1
	// description=$1
	// title=$1, description=$2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.client.Exec(ctx, query, args...)
	return err

}
