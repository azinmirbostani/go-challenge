package device

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorInvalidDeviceData   = "HTTP 400 Bad Request"
	ErrorDeviceAlreadyExists = "HTTP 400 Bad Request - Device Already Exist"
	ErrorNotFound            = "HTTP 404 Not Found"
	ErrorInternalServer      = "HTTP 500 Internal Server Error"
)

// request struct
type Device struct {
	ID          string `json:"id,omitempty"`
	DeviceModel string `json:"deviceModel,omitempty"`
	Name        string `json:"name,omitempty"`
	Note        string `json:"note,omitempty"`
	Serial      string `json:"serial,omitempty"`
}

// checking request fields
func CheckMissingFields(d Device) error {
	Flag := false
	ErrorText := "HTTP 400 Bad Request - Missing fields:"

	if d.ID == "" {
		Flag = true
		ErrorText += " 'id' "
	}
	if d.DeviceModel == "" {
		Flag = true
		ErrorText += " 'deviceModel' "
	}
	if d.Name == "" {
		Flag = true
		ErrorText += " 'name' "
	}
	if d.Note == "" {
		Flag = true
		ErrorText += " 'note' "
	}
	if d.Serial == "" {
		Flag = true
		ErrorText += " 'serial' "
	}

	if Flag == true {
		return errors.New(ErrorText)
	}

	return nil
}

// fetch device from dynamo - GetItem
func FetchDevice(ID, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Device, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(ID),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorNotFound)
	}

	if len(result.Item) == 0 {
		return nil, errors.New(ErrorNotFound)
	}

	item := new(Device)
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		return nil, errors.New(ErrorInternalServer)
	}
	return item, nil
}

// create device dynamo db - PutItem
func CreateDevice(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*Device,
	error,
) {
	var d Device

	if err := json.Unmarshal([]byte(req.Body), &d); err != nil {
		return nil, errors.New(ErrorInvalidDeviceData)
	}

	if err := CheckMissingFields(d); err != nil {
		return nil, errors.New(err.Error())
	}

	currentDevice, _ := FetchDevice(d.ID, tableName, dynaClient)
	if currentDevice != nil && len(currentDevice.ID) != 0 {
		return nil, errors.New(ErrorDeviceAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(d)

	if err != nil {
		return nil, errors.New(ErrorInternalServer)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorInternalServer)
	}
	return &d, nil
}
