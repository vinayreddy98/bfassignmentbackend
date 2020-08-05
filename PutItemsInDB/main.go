package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type ExchangeRates struct {
	Currency string
	Data     map[string]interface{}
}

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

func FetchApi() {

	//Reading data using Api
	response, err := http.Get("https://blockchain.info/ticker")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject map[string]interface{}
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		log.Fatal(err)
	}

	var write []*dynamodb.WriteRequest

	for key, value := range responseObject {

		LatestRate := ExchangeRates{Currency: key, Data: value.(map[string]interface{})}

		PutItem, _ := dynamodbattribute.MarshalMap(LatestRate)

		//doing batch write item for our input
		write = append(write, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: PutItem,
			},
		})

	}
	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			"Bitcoin": write,
		},
	}
	_, err = db.BatchWriteItem(input)
	if err != nil {
		fmt.Println(err.Error())
	}

}
func main() {
	lambda.Start(FetchApi)
}
