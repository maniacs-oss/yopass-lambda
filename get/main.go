package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/yopass/yopass-lambda"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//GetSecret returns secrets from dynamo
func GetSecret(request events.APIGatewayProxyRequest, db yopass.Database) (events.APIGatewayProxyResponse, error) {
	if _, ok := request.PathParameters["secretID"]; !ok {
		return events.APIGatewayProxyResponse{
			Body:       "SecretId missing",
			StatusCode: http.StatusBadRequest,
		}, nil
	}
	secret, err := db.Get(request.PathParameters["secretID"])
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			Body:       "Secret not found",
			StatusCode: http.StatusBadRequest,
		}, nil
	}
	resp, _ := json.Marshal(map[string]string{"secret": string(secret), "message": "OK"})
	return events.APIGatewayProxyResponse{
		Body:       string(resp),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return GetSecret(request, yopass.NewDynamo(os.Getenv("TABLE_NAME")))
	})
}
