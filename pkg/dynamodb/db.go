package db

import (
	"fmt"
	"github.com/Alwandy/system-design/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"log"
	"os"
)

type DB struct {}

var (
	creds = credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")
	sess, _ = session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_DEFAULT_REGION")),
		Credentials: creds,
	})
	tableName = "ShortUrls"
)

func (d *DB) conn() *dynamodb.DynamoDB {
	return dynamodb.New(sess)
}

func CreateTables() {
	var db = DB{}
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ShortenUrl"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("URL"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ShortenUrl"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("URL"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := db.conn().CreateTable(input)
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}
}

func CreateItem(url models.Url) error{
	var db = DB{}
	item := models.Url{
		ShortenUrl:   url.ShortenUrl,
		Url:	url.Url,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		if err != nil {
			return errors.New(fmt.Sprintf("Got error marshalling:%s", err))
		}	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = db.conn().PutItem(input)
	if err != nil {
		return errors.New(fmt.Sprintf("Got error calling PutItem:%s", err))
	}

	return nil
}

func GetItem(id string) (models.Url, error){
	var db = DB{}
	var url = fmt.Sprintf("https://bit.ly/%s", id)
	urlObject := models.Url{}
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"ShortenUrl": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(url),
					},
				},
			},
		},
	}

	var resp, err = db.conn().Query(queryInput)
	if err != nil {
		return urlObject, errors.New(fmt.Sprintf("Got error calling GetItem:%s", err))
	}

	// Deference the pointer
	urlObject.ShortenUrl = *resp.Items[0]["ShortenUrl"].S
	urlObject.Url = *resp.Items[0]["URL"].S

	return urlObject, nil
}