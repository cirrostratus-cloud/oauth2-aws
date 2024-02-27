package main

import (
	"context"
	"encoding/json"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/cirrostratus-cloud/common/event"
	"github.com/cirrostratus-cloud/oauth2-aws/user/repository"
	user_service "github.com/cirrostratus-cloud/oauth2-aws/user/service"
	"github.com/cirrostratus-cloud/oauth2/user"
	log "github.com/sirupsen/logrus"
)

var re *regexp.Regexp = regexp.MustCompile(`arn:aws:sqs:[a-z]{2}-[a-z]*-[0-9]:[0-9]{12}:cirrostratus-oauth2_user_(.*)`)

var snsEventBus *user_service.SNSEventBus

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
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("Error loading AWS config")
		panic(err)
	}
	topicArnPrefix := os.Getenv("TOPIC_ARN_PREFIX")
	if topicArnPrefix == "" {
		log.Fatal("TOPIC_ARN_PREFIX is required")
		panic("TOPIC_ARN_PREFIX is required")
	}
	emailFrom := os.Getenv("SES_EMAIL_FROM")
	if emailFrom == "" {
		log.Fatal("SES_EMAIL_FROM is required")
		panic("SES_EMAIL_FROM is required")
	}
	snsClient := sns.NewFromConfig(cfg)
	sesClient := ses.NewFromConfig(cfg)
	dynamodbClient := dynamodb.NewFromConfig(cfg)
	snsEventBus = user_service.NewSNSEventBus(snsClient, topicArnPrefix)
	userRepository := repository.NewDynamoUserRepository(dynamodbClient)
	mailService := user_service.NewSESMailService(sesClient)
	user.NewNotifyPasswordChangedService(userRepository, mailService, snsEventBus, emailFrom)
	user.NewNotifyPasswordRecoveredService(userRepository, mailService, snsEventBus, emailFrom)
	user.NewNotifyUserCreatedService(userRepository, mailService, snsEventBus, emailFrom)
}

func handler(ctx context.Context, req events.SQSEvent) {
	for _, record := range req.Records {
		log.WithFields(log.Fields{
			"MessageID":      record.MessageId,
			"Body":           record.Body,
			"EventSourceARN": record.EventSourceARN,
		}).Info("Received message")
		var recordBody map[string]interface{}
		err := json.Unmarshal([]byte(record.Body), &recordBody)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error unmarshalling record body")
			return
		}
		queueArn := record.EventSourceARN
		matches := re.FindStringSubmatch(queueArn)
		if matches == nil {
			log.WithFields(log.Fields{
				"QueueArn": queueArn,
			}).Error("Regex did not match queueArn")
			return
		}
		if len(matches) < 2 {
			log.WithFields(log.Fields{
				"QueueArn": queueArn,
			}).Error("Error parsing queueArn")
			return
		}
		eventName := "user/" + matches[1]
		log.WithFields(log.Fields{
			"EventName": eventName,
		}).Info("Triggering event")
		err = snsEventBus.Trigger(event.EventName(eventName), recordBody["Message"].(string))
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Error triggering event")
		}
	}
}

func main() {
	lambda.Start(handler)
}
