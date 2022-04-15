package device

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (self *MockDynamoDB) GetItem(input *dynamodb.GetItemInput) (output *dynamodb.GetItemOutput, err error) {

	ID := input.Key["id"].S
	ItemOutput := dynamodb.GetItemOutput{}

	if *ID == "id1" {
		ItemOutput.SetItem(
			map[string]*dynamodb.AttributeValue{
				"id":          &dynamodb.AttributeValue{S: ID},
				"deviceModel": &dynamodb.AttributeValue{S: aws.String("deviceModel1")},
				"name":        &dynamodb.AttributeValue{S: aws.String("Sensor")},
				"note":        &dynamodb.AttributeValue{S: aws.String("Testing a sensor.")},
				"serial":      &dynamodb.AttributeValue{S: aws.String("A020000102")},
			},
		)
	}
	return &ItemOutput, err
}

func (self *MockDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	MockOutput := new(dynamodb.PutItemOutput)
	return MockOutput, nil
}

func TestFetchDevice(t *testing.T) {
	type args struct {
		ID         string
		tableName  string
		dynaClient dynamodbiface.DynamoDBAPI
	}
	tests := []struct {
		name    string
		args    args
		want    *Device
		wantErr bool
	}{
		{
			name:    "* Test: Missing Fields: ID *",
			args:    args{ID: "", tableName: "go-challenge-devices", dynaClient: &MockDynamoDB{}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "* Test: Not Found ID *",
			args:    args{ID: "id2", tableName: "go-challenge-devices", dynaClient: &MockDynamoDB{}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "* Test: Find ID Correctly *",
			args: args{ID: "id1", tableName: "go-challenge-devices", dynaClient: &MockDynamoDB{}},
			want: &Device{
				ID:          "id1",
				DeviceModel: "deviceModel1",
				Name:        "Sensor",
				Note:        "Testing a sensor.",
				Serial:      "A020000102",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchDevice(tt.args.ID, tt.args.tableName, tt.args.dynaClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateDevice(t *testing.T) {
	type args struct {
		req        events.APIGatewayProxyRequest
		tableName  string
		dynaClient dynamodbiface.DynamoDBAPI
	}
	tests := []struct {
		name    string
		args    args
		want    *Device
		wantErr bool
	}{
		{
			name:    "* Test: Missing Fields: ID *",
			args:    args{req: events.APIGatewayProxyRequest{Body: "{\"id\":\"\", \"deviceModel\":\"deviceModel1\", \"name\":\"Sensor\", \"note\":\"Testing a sensor.\", \"serial\":\"A020000102\"}"}, tableName: "go-challenge-devices", dynaClient: &MockDynamoDB{}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "* Test: Missing Fields: deviceModel *",
			args:    args{req: events.APIGatewayProxyRequest{Body: "{\"id\":\"id1\", \"deviceModel\":\"\", \"name\":\"Sensor\", \"note\":\"Testing a sensor.\", \"serial\":\"A020000102\"}"}, tableName: "go-challenge-devices", dynaClient: &MockDynamoDB{}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "* Test: Missing Fields: Name *",
			args:    args{req: events.APIGatewayProxyRequest{Body: "{\"id\":\"id1\", \"deviceModel\":\"deviceModel1\", \"name\":\"\", \"note\":\"Testing a sensor.\", \"serial\":\"A020000102\"}"}, tableName: "go-challenge-devices", dynaClient: &MockDynamoDB{}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "* Test: Missing Fields: Note *",
			args:    args{req: events.APIGatewayProxyRequest{Body: "{\"id\":\"id1\", \"deviceModel\":\"deviceModel1\", \"name\":\"Sensor\", \"note\":\"\", \"serial\":\"A020000102\"}"}, tableName: "go-challenge-devices", dynaClient: &MockDynamoDB{}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "* Test: Missing Fields: Serial *",
			args:    args{req: events.APIGatewayProxyRequest{Body: "{\"id\":\"id1\", \"deviceModel\":\"deviceModel1\", \"name\":\"Sensor\", \"note\":\"Testing a sensor.\", \"serial\":\"\"}"}, tableName: "go-challenge-devices", dynaClient: &MockDynamoDB{}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "* Test: Device Add Correctly *",
			args: args{req: events.APIGatewayProxyRequest{Body: "{\"id\":\"id2\", \"deviceModel\":\"deviceModel2\", \"name\":\"Sensor\", \"note\":\"Testing a sensor.\", \"serial\":\"A020000102\"}"}, tableName: "go-challenge-devices", dynaClient: &MockDynamoDB{}},
			want: &Device{
				ID:          "id2",
				DeviceModel: "deviceModel2",
				Name:        "Sensor",
				Note:        "Testing a sensor.",
				Serial:      "A020000102",
			},
			wantErr: false,
		},
		{
			name:    "* Test: Device Already Exist *",
			args:    args{req: events.APIGatewayProxyRequest{Body: "{\"id\":\"id1\", \"deviceMode1\":\"deviceModel1\", \"name\":\"Sensor\", \"note\":\"Testing a sensor.\", \"serial\":\"A020000102\"}"}, tableName: "go-challenge-devices", dynaClient: &MockDynamoDB{}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDevice(tt.args.req, tt.args.tableName, tt.args.dynaClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}
