package api

import (
	"fmt"
	"lamda-func/database"
	"lamda-func/types"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error {
	if event.Username == "" || event.Password == "" {
		return fmt.Errorf("Invalid request: Username or password is empty")
	}

	userExists, err := api.dbStore.DoesUserExist(event.Username)
	if err != nil {
		return fmt.Errorf("Error: %w", err)
	}

	if userExists {
		return fmt.Errorf("User with that username already exist")
	}

	err = api.dbStore.InsertUser(event)
	if err != nil {
		return fmt.Errorf("Error creating user: %w", err)
	}

	return nil
}
