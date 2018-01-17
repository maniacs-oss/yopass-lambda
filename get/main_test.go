package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

type mockDB struct{}

func (db *mockDB) Get(key string) (string, error) {
	return "wop", nil
}
func (db *mockDB) Put(key, value string, expiration int32) error {
	return nil
}
func (db *mockDB) Delete(key string) error {
	return nil
}

func TestGetSecret(t *testing.T) {
	params := map[string]string{"secretID": "36f6d7df-348b-4b38-9fbc-8b193a781d18"}
	db := &mockDB{}
	resp, err := GetSecret(events.APIGatewayProxyRequest{
		Body:           "",
		PathParameters: params,
	}, db)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp.Body)
}
