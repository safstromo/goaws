package api

import (
	"encoding/json"
	"fmt"
	"lamda-func/database"
	"lamda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser types.RegisterUser

	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	// TODO: Cleanup refactor
	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request: Username or password is empty",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	userExists, err := api.dbStore.DoesUserExist(registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if userExists {
		return events.APIGatewayProxyResponse{
			Body:       "User with that username already exist",
			StatusCode: http.StatusConflict,
		}, nil
	}

	user, err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("Could not create new user %w", err)
	}

	err = api.dbStore.InsertUser(*user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("Error inserting user: %w", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       "Successfully registered user",
		StatusCode: http.StatusCreated,
	}, nil
}

func (api ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := api.dbStore.GetUser(loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if !types.ValidatePassword(user.PasswordHash, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body:       "Invalid credentials",
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "Successfully logged in",
		StatusCode: http.StatusOK,
	}, nil
}
