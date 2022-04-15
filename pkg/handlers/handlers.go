package handlers

import (
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/azinmirbostani/go-challenge/pkg/device"
)

var ErrorMethodNotAllowed = "Method Not Allowed"
var ErrorNotFound = "HTTP 404 Not Found"

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

// find device - call device fetch dynamo db
func GetDevice(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse, error,
) {
	ID := req.PathParameters["id"]
	if ID == "" {
		return nil, errors.New(ErrorNotFound)
	}

	_ID := "/devices/" + ID
	result, err := device.FetchDevice(_ID, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}
	return apiResponse(http.StatusOK, result)
}

// create device - call create device dynamo db
func CreateDevice(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse, error,
) {
	result, err := device.CreateDevice(req, tableName, dynaClient)
	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{
			aws.String(err.Error()),
		})
	}
	return apiResponse(http.StatusCreated, result)
}

// handle not defined methods
func UnhandledMethod() (*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
