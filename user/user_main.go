package main

import (
	"context"
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
	stage := os.Getenv("AWS_STAGE")
	app = fiber.New()
	setUp(app, stage)
	fiberLambda = fiberadapter.New(app)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.
		WithField("Path", req.RequestContext.Path).
		Info("Processing request.")
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
