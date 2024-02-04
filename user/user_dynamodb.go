package main

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cirrostratus-cloud/oauth2/user"
)

var tableName string = fmt.Sprintf("%s-%s", os.Getenv("CIRROSTRATUS_OAUTH2_MODULE_NAME"), os.Getenv("CIRROSTRATUS_OUTH2_USER_TABLE"))
var getOneScanLimit int32 = 1

type dynamoUserRepository struct {
	client *dynamodb.Client
}

func newDynamoUserRepository(client *dynamodb.Client) *dynamoUserRepository {
	return &dynamoUserRepository{client: client}
}

func (u *dynamoUserRepository) CreateUser(user user.User) (user.User, error) {
	_, err := u.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &tableName,
		Item: map[string]types.AttributeValue{
			"id":        &types.AttributeValueMemberS{Value: user.GetID()},
			"email":     &types.AttributeValueMemberS{Value: user.GetEmail()},
			"password":  &types.AttributeValueMemberS{Value: user.GetPassword()},
			"enabled":   &types.AttributeValueMemberBOOL{Value: user.IsEnabled()},
			"firstName": &types.AttributeValueMemberS{Value: user.GetFirstName()},
			"lastName":  &types.AttributeValueMemberS{Value: user.GetLastName()},
		},
	})
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *dynamoUserRepository) GetUserByID(userID string) (user.User, error) {
	output, err := u.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return user.User{}, err
	}
	if len(output.Item) == 0 {
		return user.User{}, fmt.Errorf("user with id %s not found", userID)
	}
	userEntity, err := user.NewUser(
		output.Item["id"].(*types.AttributeValueMemberS).Value,
		output.Item["email"].(*types.AttributeValueMemberS).Value,
		output.Item["password"].(*types.AttributeValueMemberS).Value,
	)
	enabled := output.Item["enabled"].(*types.AttributeValueMemberBOOL).Value
	if enabled {
		userEntity.EnableUser()
	} else {
		userEntity.DisableUser()
	}
	firstName := output.Item["firstName"].(*types.AttributeValueMemberS).Value
	lastName := output.Item["lastName"].(*types.AttributeValueMemberS).Value
	userEntity.UpdateUserProfile(firstName, lastName)
	if err != nil {
		return user.User{}, err
	}
	return userEntity, nil
}
func (u *dynamoUserRepository) UpdateUser(user user.User) (user.User, error) {
	return user, nil
}
func (u *dynamoUserRepository) GetUserByEmail(email string) (user.User, error) {
	filterExpression := expression.Name("email").Equal(expression.Value(email))
	exp, err := expression.NewBuilder().WithFilter(filterExpression).Build()
	if err != nil {
		return user.User{}, err
	}
	output, err := u.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:                 &tableName,
		FilterExpression:          exp.Filter(),
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
		Limit:                     &getOneScanLimit,
	})
	if err != nil {
		log.
			WithField("email", email).
			Error(err)
		return user.User{}, err
	}
	if len(output.Items) == 0 {
		return user.User{}, nil
	}
	if len(output.Items) > 1 {
		err = fmt.Errorf("more than one user found with email %s", email)
		log.
			WithFields(log.Fields{
				"email": email,
				"count": len(output.Items),
			}).
			Error(err)
		return user.User{}, err
	}
	foundedUser, err := user.NewUser(output.Items[0]["id"].(*types.AttributeValueMemberS).Value, output.Items[0]["email"].(*types.AttributeValueMemberS).Value, output.Items[0]["password"].(*types.AttributeValueMemberS).Value)
	if err != nil {
		return user.User{}, err
	}
	enabled := output.Items[0]["enabled"].(*types.AttributeValueMemberBOOL).Value
	if enabled {
		foundedUser.EnableUser()
	} else {
		foundedUser.DisableUser()
	}
	firstName := output.Items[0]["firstName"].(*types.AttributeValueMemberS).Value
	lastName := output.Items[0]["lastName"].(*types.AttributeValueMemberS).Value
	foundedUser.UpdateUserProfile(firstName, lastName)
	return foundedUser, nil
}

func (u *dynamoUserRepository) ExistUserByEmail(email string) (bool, error) {
	filterExpression := expression.Name("email").Equal(expression.Value(email))
	exp, err := expression.NewBuilder().WithFilter(filterExpression).Build()
	if err != nil {
		return false, err
	}
	output, err := u.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:                 &tableName,
		FilterExpression:          exp.Filter(),
		ExpressionAttributeNames:  exp.Names(),
		ExpressionAttributeValues: exp.Values(),
	})
	if err != nil {
		log.
			WithField("email", email).
			Error(err)
		return false, err
	}
	if len(output.Items) == 0 {
		return false, nil
	}
	return true, nil
}
