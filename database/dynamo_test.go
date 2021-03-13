package database

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type stubDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func getAttrValue() map[string]*dynamodb.AttributeValue {
	// Initialize response
	key := dynamodb.AttributeValue{}
	key.SetS("shortId")
	url := dynamodb.AttributeValue{}
	url.SetS("www.sample-long-url-address-for-test.com")
	resp := make(map[string]*dynamodb.AttributeValue)
	resp["id"] = &key
	resp["URL"] = &url
	return resp
}

func (m *stubDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	resp := getAttrValue()

	output := &dynamodb.GetItemOutput{
		Item: resp,
	}

	return output, nil
}

func (m *stubDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	resp := getAttrValue()

	output := &dynamodb.PutItemOutput{
		Attributes: resp,
	}

	return output, nil
}

func TestGet(t *testing.T) {
	stub := &stubDynamoDB{}

	res, err := getItem(stub, "table", "id")
	if err != nil {
		t.Errorf("Error calling Dynamodb %d", err)
	}
	keyRes := *res.Item["id"].S
	valRes := *res.Item["URL"].S
	if keyRes != "shortId" {
		t.Errorf("Wrong key returned. Shoule be id, was %s", keyRes)
	}
	if valRes != "www.sample-long-url-address-for-test.com" {
		t.Errorf("Wrong value returned. Shoule be URL, was %s", valRes)
	}
}

func TestPost(t *testing.T) {
	stub := &stubDynamoDB{}
	attrVal := getAttrValue()
	res, _ := saveItem(stub, "table", attrVal)
	keyRes := *res.Attributes["id"].S
	valRes := *res.Attributes["URL"].S

	if keyRes != "shortId" {
		t.Errorf("Wrong key returned. Shoule be id, was %s", keyRes)
	}
	if valRes != "www.sample-long-url-address-for-test.com" {
		t.Errorf("Wrong value returned. Shoule be URL, was %s", valRes)
	}
}
