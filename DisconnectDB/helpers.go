package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type ConnectionItem struct {
	ConnectionID string `json:"connectionID"`
}

const (
	AccessKeyID        = "AKIAII5ZHONLGXRN7B7A"
	SecretAccessKey    = "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"
	APIGatewayEndpoint = "https://u3fen3k27a.execute-api.us-east-1.amazonaws.com/Final"
	Region             = "us-east-1"
)

func NewDynamoDBSession() *dynamodb.DynamoDB {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(Region),
		Credentials: credentials.NewStaticCredentials(AccessKeyID, SecretAccessKey, ""),
	})
	return dynamodb.New(sess)
}

func NewAPIGatewaySession() *apigatewaymanagementapi.ApiGatewayManagementApi {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(Region),
		Credentials: credentials.NewStaticCredentials(AccessKeyID, SecretAccessKey, ""),
		Endpoint:    aws.String(APIGatewayEndpoint),
	})
	return apigatewaymanagementapi.New(sess)
}
