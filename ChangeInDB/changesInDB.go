package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func main() {
	lambda.Start(HandleMessage)
}

type DynamoEventChange struct {
	NewImage map[string]*dynamodb.AttributeValue `json:"NewImage"`
}

type DynamoEventRecord struct {
	Change    DynamoEventChange `json:"dynamodb"`
	EventName string            `json:"eventName"`
	EventID   string            `json:"eventID"`
}

type DynamoEvent struct {
	Records []DynamoEventRecord `json:"records"`
}

type ConnectionItem struct {
	ConnectionID string `json:"connectionID"`
}

type ExchangeRates struct {
	Currency string
	Data     map[string]interface{}
}

type MessageData struct {
	Message []ExchangeRates `json:"message"`
}

func HandleMessage(req DynamoEvent) (events.APIGatewayProxyResponse, error) {

	dynamodbSession := NewDynamoDBSession()

	// Read the chat-connections table
	connectionInputs := &dynamodb.ScanInput{
		TableName: aws.String("connections"),
	}

	bitcoinInput := &dynamodb.ScanInput{
		TableName: aws.String("Bitcoin"),
	}
	scanconnInputs, _ := dynamodbSession.Scan(connectionInputs)

	scanBitcoinInputs, _ := dynamodbSession.Scan(bitcoinInput)

	// Parse the table data in the output variable
	var output []ConnectionItem
	dynamodbattribute.UnmarshalListOfMaps(scanconnInputs.Items, &output)

	var responseObject []ExchangeRates
	dynamodbattribute.UnmarshalListOfMaps(scanBitcoinInputs.Items, &responseObject)

	apigatewaySession := NewAPIGatewaySession()

	// Encode the message data with details obtained from dynamodbtable
	messageData := &MessageData{}
	for _, rcd := range req.Records {
		rdata := &ExchangeRates{}
		dynamodbattribute.UnmarshalMap(rcd.Change.NewImage, rdata)
		messageData.Message = append(messageData.Message, *rdata)
	}

	jsonData, _ := json.Marshal((messageData.Message))

	// Send the message for each connection ID
	for _, item := range output {
		connectionInput := &apigatewaymanagementapi.PostToConnectionInput{
			ConnectionId: aws.String(item.ConnectionID),
			Data:         jsonData,
		}

		_, err := apigatewaySession.PostToConnection(connectionInput)
		if err != nil {
			fmt.Println(err)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}
