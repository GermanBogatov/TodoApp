package test

import (
	"errors"
	"github.com/GermanBogatov/TodoApp/app/internal/handler"
	"github.com/GermanBogatov/TodoApp/app/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidate(t *testing.T) {
	var user = model.UserDTO{
		Name:     "German",
		Username: "Bogatov",
		Password: "23Bogatov",
	}

	err := handler.ValidateRequestSign(user)
	require.NoError(t, err)
}

func TestValidateError(t *testing.T) {
	cases := []struct {
		name   string
		user   model.UserDTO
		expErr error
	}{
		{
			name: "bad id test",
			user: model.UserDTO{
				Id:       12,
				Name:     "German",
				Username: "Bogatov",
				Password: "23Bogatov",
			},
			expErr: errors.New("Invalid id in request!"),
		},
		{
			name: "bad username test",
			user: model.UserDTO{
				Id:       0,
				Name:     "German",
				Username: "",
				Password: "23Bogatov",
			},
			expErr: errors.New("Invalid Username!"),
		},
		{
			name: "bad name test",
			user: model.UserDTO{
				Id:       0,
				Name:     "",
				Username: "Bogatov",
				Password: "23Bogatov",
			},
			expErr: errors.New("Invalid Name!"),
		},
		{
			name: "bad password test",
			user: model.UserDTO{
				Id:       0,
				Name:     "German",
				Username: "Bogatov",
				Password: "",
			},
			expErr: errors.New("Invalid Password!"),
		},
		{
			name: "bad length password test",
			user: model.UserDTO{
				Id:       0,
				Name:     "German",
				Username: "Bogatov",
				Password: "Bad",
			},
			expErr: errors.New("Invalid length Password!"),
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			err := handler.ValidateRequestSign(v.user)
			require.Error(t, err)
			require.EqualError(t, v.expErr, err.Error())
		})
	}
}
