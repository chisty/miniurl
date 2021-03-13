package database

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/chisty/shortlink/model"
)

type dymamoDB struct {
	table   string
	context dynamodbiface.DynamoDBAPI
	logger  *log.Logger
}

//NewDynamoDB ---
func NewDynamoDB(tbl string, log *log.Logger, rgn, acsKey, secKey string) DB {
	return &dymamoDB{
		table:   tbl,
		logger:  log,
		context: createDBClient(rgn, acsKey, secKey),
	}
}

func (db *dymamoDB) Get(id string) (*model.ShortLink, error) {
	db.logger.Println("get from ddb: ", id)

	item, err := getItem(db.context, db.table, id)
	if err != nil {
		return nil, err
	}

	slink := model.ShortLink{}
	err = dynamodbattribute.UnmarshalMap(item.Item, &slink)
	if err != nil {
		db.logger.Fatal(err)
		return nil, err
	}

	db.logger.Printf("get from ddb success with value: %s", slink.URL)
	return &slink, nil
}

func (db *dymamoDB) Save(data *model.ShortLink) error {
	db.logger.Println("save in ddb: ", data.ID)

	attrVal, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return err
	}

	_, err = saveItem(db.context, db.table, attrVal)
	if err != nil {
		return err
	}

	db.logger.Printf("data saved with id: %s\n", data.ID)

	return nil
}

func getItem(dbApi dynamodbiface.DynamoDBAPI, tbl string, id string) (*dynamodb.GetItemOutput, error) {
	item, err := dbApi.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tbl),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return item, nil
}

func saveItem(dbApi dynamodbiface.DynamoDBAPI, tbl string, attrVal map[string]*dynamodb.AttributeValue) (*dynamodb.PutItemOutput, error) {
	item, err := dbApi.PutItem(&dynamodb.PutItemInput{
		Item:      attrVal,
		TableName: aws.String(tbl),
	})

	if err != nil {
		return nil, err
	}

	return item, nil
}

func createDBClient(region, accessKey, secretKey string) *dynamodb.DynamoDB {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}))

	return dynamodb.New(sess)
}
