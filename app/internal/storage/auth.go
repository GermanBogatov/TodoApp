package storage

import (
	"context"
	"fmt"
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
	"github.com/GermanBogatov/TodoApp/app/pkg/postgresql"
	"github.com/jackc/pgconn"
)

type repositoryAuth struct {
	client postgresql.ClientPostgres
	logger logging.Logger
}

func NewAuthorization(client postgresql.ClientPostgres, logger logging.Logger) *repositoryAuth {
	return &repositoryAuth{
		client: client,
		logger: logger,
	}
}

func (r *repositoryAuth) CreateUser(ctx context.Context, user model.User) (int, error) {

	q := `
    INSERT INTO users 
    	(name,username,password_hash) 
    VALUES 
		($1,$2,$3) 
    RETURNING id
		`

	if err := r.client.QueryRow(ctx, q, user.Name, user.Username, user.Password).Scan(&user.Id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return 0, newErr
		}
		return 0, err
	}
	return user.Id, nil
}

func (r *repositoryAuth) GetUser(ctx context.Context, username, password string) (model.User, error) {
	var user model.User

	q := `
	SELECT id, username
	FROM users
	WHERE username=$1 AND password_hash=$2
		`

	err := r.client.QueryRow(ctx, q, username, password).Scan(&user.Id, &user.Username)
	if err != nil {
		return model.User{}, err
	}

	return user, nil

}
