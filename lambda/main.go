package main

import (
	"fmt"
	"lamda-func/app"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

// Take in payload and something with it
func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	return fmt.Sprintf("Successfully called by - %s", event.Username), nil
}

// Build command
// GOOS=linux GOARCH=amd64 go build -o bootstrap
// zip function.zip bootstrap

func main() {
	myApp := app.NewApp()
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return myApp.ApiHandler.RegisterUserHandler(request)
		case "/login":
			return myApp.ApiHandler.LoginUser(request)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "Not Found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	})
}
