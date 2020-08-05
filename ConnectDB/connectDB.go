package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	lambda.Start(HandleConnect)
}

func HandleConnect(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	connectionItem := ConnectionItem{
		ConnectionID: request.RequestContext.ConnectionID,
	}
	attributeValues, _ := dynamodbattribute.MarshalMap(connectionItem)

	input := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: aws.String("connections"),
	}

	dynamodbSession := NewDynamoDBSession()
	dynamodbSession.PutItem(input)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}
