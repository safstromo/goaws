package database

import (
	"fmt"
	"lamda-func/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	TABLE_NAME = "userTable"
)

type DynamoDBClient struct {
	databaseStore *dynamodb.DynamoDB
}

// TODO: Refactor interface to own file
type UserStore interface {
	DoesUserExist(username string) (bool, error)
	GetUser(username string) (types.User, error)
	InsertUser(user types.User) error
}

func NewDynamoDBClient() DynamoDBClient {
	dbSession := session.Must(session.NewSession())

	db := dynamodb.New(dbSession)

	return DynamoDBClient{
		databaseStore: db,
	}
}

func (u DynamoDBClient) DoesUserExist(username string) (bool, error) {
	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return true, err
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (u DynamoDBClient) InsertUser(user types.User) error {
	// assemble item
	// TODO: encrypt
	item := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(user.Username),
			},
			"password": {
				S: aws.String(user.PasswordHash),
			},
		},
	}

	_, err := u.databaseStore.PutItem(item)
	if err != nil {
		return err
	}

	return nil
}

func (u DynamoDBClient) GetUser(username string) (types.User, error) {
	var user types.User

	result, err := u.databaseStore.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	})
	if err != nil {
		return user, err
	}

	if result.Item == nil {
		return user, fmt.Errorf("User does not exist")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}
