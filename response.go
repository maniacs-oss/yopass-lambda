package yopass

import "github.com/aws/aws-lambda-go/events"

// Response function that returns reponse with correct headers
func Response(message string, statusCode int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       message,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		StatusCode: statusCode,
	}, nil
}
