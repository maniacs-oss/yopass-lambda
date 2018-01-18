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
		return yopass.Response("SecretId missing", http.StatusBadRequest)
	}
	secret, err := db.Get(request.PathParameters["secretID"])
	if err != nil {
		fmt.Println(err)
		return yopass.Response("Secret not found", http.StatusBadRequest)
	}
	resp, _ := json.Marshal(map[string]string{"secret": string(secret), "message": "OK"})
	return yopass.Response(string(resp), 200)
}

func main() {
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return GetSecret(request, yopass.NewDynamo(os.Getenv("TABLE_NAME")))
	})
}
