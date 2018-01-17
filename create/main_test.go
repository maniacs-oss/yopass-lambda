package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

type mockDB struct{}

func (db *mockDB) Get(key string) (string, error) {
	return "", nil
}
func (db *mockDB) Put(key, value string, expiration int32) error {
	return nil
}
func (db *mockDB) Delete(key string) error {
	return nil
}
func TestCreateSecret(t *testing.T) {
	tt := []struct {
		name       string
		statusCode int
		body       string
		output     string
	}{
		{
			name:       "validRequest",
			statusCode: 200,
			body:       `{"secret": "hello world", "expiration": 3600}`,
			output:     "OK",
		},
		{
			name:       "invalid json",
			statusCode: 400,
			body:       `{fooo`,
			output:     "Unable to parse json",
		},
		{
			name:       "message too long",
			statusCode: 400,
			body:       `{"expiration": 3600, "secret": "` + strings.Join(make([]string, 12000), "x") + `"}`,
			output:     "Message is too long",
		},
	}

	for _, tc := range tt {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			request := events.APIGatewayProxyRequest{Body: tc.body}
			db := &mockDB{}
			resp, _ := CreateSecret(request, db)
			var response struct {
				Message string `json:"message"`
				Key     string `json:"key"`
			}
			json.Unmarshal([]byte(resp.Body), &response)
			if response.Message != tc.output {
				t.Fatalf(`Expected body "%s"; got "%s"`, tc.output, response.Message)
			}
			if resp.StatusCode != tc.statusCode {
				t.Fatalf(`Expected status code %d; got "%d"`, tc.statusCode, resp.StatusCode)
			}
		})
	}
}
