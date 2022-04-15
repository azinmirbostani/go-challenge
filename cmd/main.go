package main

import (
	"os"

	"github.com/azinmirbostani/go-challenge/pkg/handlers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)})

	if err != nil {
		return
	}

	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

const tableName = "go-challenge-devices"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// checking request
	switch req.HTTPMethod {
	case "GET":
		return handlers.GetDevice(req, tableName, dynaClient)
	case "POST":
		return handlers.CreateDevice(req, tableName, dynaClient)
	default:
		return handlers.UnhandledMethod()
	}
}
