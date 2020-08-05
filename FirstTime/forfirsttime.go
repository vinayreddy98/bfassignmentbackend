package main

import (
	"context"
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
	lambda.Start(HandlefirsttimeConnection)
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

func HandlefirsttimeConnection(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Parse the request

	dynamodbSession := NewDynamoDBSession()

	// Read the connections table
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

	jsonData, _ := json.Marshal(responseObject)

	// Send the message for the  connected connection ID
	for _, item := range output {
		if item.ConnectionID == request.RequestContext.ConnectionID {

			connectionInput := &apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: aws.String(item.ConnectionID),
				Data:         jsonData,
			}

			_, err := apigatewaySession.PostToConnection(connectionInput)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			continue
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
	}, nil
}
