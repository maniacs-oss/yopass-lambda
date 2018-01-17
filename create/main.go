package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	uuid "github.com/satori/go.uuid"
	"github.com/yopass/yopass-lambda"
)

// validExpiration validates that expiration is either
// 3600(1hour), 86400(1day) or 604800(1week)
func validExpiration(expiration int32) bool {
	for _, ttl := range []int32{3600, 86400, 604800} {
		if ttl == expiration {
			return true
		}
	}
	return false
}

// CreateSecret creates secret
func CreateSecret(request events.APIGatewayProxyRequest, db yopass.Database) (events.APIGatewayProxyResponse, error) {
	var secret struct {
		Message    string `json:"secret"`
		Expiration int32  `json:"expiration"`
	}
	err := json.Unmarshal([]byte(request.Body), &secret)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       `{"message": "Unable to parse json"}`,
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if !validExpiration(secret.Expiration) {
		return events.APIGatewayProxyResponse{
			Body:       `{"message": "Invalid expiration specified"}`,
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	if len(secret.Message) > 10000 {
		return events.APIGatewayProxyResponse{
			Body:       `{"message": "Message is too long"}`,
			StatusCode: http.StatusBadRequest,
		}, nil
	}
	key := uuid.NewV4().String()
	err = db.Put(key, secret.Message, secret.Expiration)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			Body:       `{"message": "Failed to store secret in database"}`,
			StatusCode: http.StatusInternalServerError,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(`{"key": "%s", "message": "OK"}`, key),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return CreateSecret(request, yopass.NewDynamo(os.Getenv("TABLE_NAME")))
	})
}
