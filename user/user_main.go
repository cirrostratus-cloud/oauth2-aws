package main

import (
	"context"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	"github.com/cirrostratus-cloud/oauth2/user"
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
	inValue, err := strconv.Atoi(os.Getenv("USER_MIN_PASSWORD_LENGTH"))
	if err != nil {
		log.Fatal("Error parsing USER_MIN_PASSWORD_LENGTH")
		panic(err)
	}
	minPasswordLength := inValue
	boolValue, err := strconv.ParseBool(os.Getenv("USER_UPPER_CASE_REQUIRED"))
	if err != nil {
		log.Fatal("Error parsing USER_UPPER_CASE_REQUIRED")
		panic(err)
	}
	upperCaseRequired := boolValue
	boolValue, err = strconv.ParseBool(os.Getenv("USER_LOWER_CASE_REQUIRED"))
	if err != nil {
		log.Fatal("Error parsing USER_LOWER_CASE_REQUIRED")
		panic(err)
	}
	lowerCaseRequired := boolValue
	boolValue, err = strconv.ParseBool(os.Getenv("USER_NUMBER_REQUIRED"))
	if err != nil {
		log.Fatal("Error parsing USER_NUMBER_REQUIRED")
		panic(err)
	}
	numberRequired := boolValue
	boolValue, err = strconv.ParseBool(os.Getenv("USER_SPECIAL_CHARACTER_REQUIRED"))
	if err != nil {
		log.Fatal("Error parsing USER_SPECIAL_CHARACTER_REQUIRED")
		panic(err)
	}
	specialCharacterRequired := boolValue
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("Error loading AWS config")
		panic(err)
	}
	dynamodbClient := dynamodb.NewFromConfig(cfg)
	userRepository := newDynamoUserRepository(dynamodbClient)
	createUserService := user.NewCreateUserService(userRepository, minPasswordLength, upperCaseRequired, lowerCaseRequired, numberRequired, specialCharacterRequired)
	getUserUseCase := user.NewGetUserService(userRepository)
	updateProfileUseCase := user.NewUpdateUserProfileService(userRepository)
	userAPI := newUserAPI(createUserService, getUserUseCase,updateProfileUseCase)
	app = fiber.New()
	userAPI.setUp(app, stage)
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
