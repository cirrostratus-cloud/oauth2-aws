package main

import (
	"github.com/aws/aws-lambda-go/events"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Hello Client",
		StatusCode: 200,
	}, nil
}
