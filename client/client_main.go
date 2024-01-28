package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

var fiberLambda *fiberadapter.FiberLambda
var app *fiber.App

func init() {

	app = fiber.New()
	stage := os.Getenv("AWS_STAGE")
	setUp(app, stage)

	fiberLambda = fiberadapter.New(app)
	logLevel := os.Getenv("LOG_LEVEL")
	log.SetOutput(os.Stdout)
	switch logLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Processing request path %s.\n", req.RequestContext.Path)
	return fiberLambda.ProxyWithContext(ctx, req)
}

func main() {
	stage := os.Getenv("AWS_STAGE")
	if stage == "local" {
		log.Fatal(app.Listen(":3000"))
	} else {
		lambda.Start(Handler)
	}
}
