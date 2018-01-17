package yopass

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Database interface
type Database interface {
	Get(key string) (string, error)
	Put(key, value string, expiration int32) error
	Delete(key string) error
}

// Dynamo implementation
type Dynamo struct {
	tableName string
	svc       *dynamodb.DynamoDB
}

// NewDynamo returns a database client
func NewDynamo(tableName string) Database {
	return &Dynamo{tableName: tableName, svc: dynamodb.New(session.New())}
}

// Get item from dynamo
func (d *Dynamo) Get(key string) (string, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
		},
		TableName: aws.String(d.tableName),
	}
	result, err := d.svc.GetItem(input)
	if err != nil {
		return "", err
	}
	if len(result.Item) == 0 {
		return "", fmt.Errorf("Key not found in database")
	}
	return *result.Item["secret"].S, nil
}

// Delete item
func (d *Dynamo) Delete(key string) error {
	return nil
}

// Put item in Dynamo
func (d *Dynamo) Put(key, value string, expiration int32) error {
	input := &dynamodb.PutItemInput{
		// TABLE GENERATED NAME
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
			"secret": {
				S: aws.String(value),
			},
			"ttl": {
				N: aws.String(
					fmt.Sprintf(
						"%d", time.Now().Unix()+int64(expiration))),
			},
		},
		TableName: aws.String(d.tableName),
	}
	_, err := d.svc.PutItem(input)
	return err
}
