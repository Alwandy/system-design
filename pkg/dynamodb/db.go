package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/Alwandy/system-design/pkg/dynamodb"
	"log"
)

var (
	creds = credentials.NewStaticCredentials("", "", "")
	sess, _ = session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: creds,
	})
)

func (d *db) conn() *dynamodb.DynamoDB {
	return dynamodb.New(sess)
}

func (d *db) createTables() {
	tableName := "ShortUrls"
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

	_, err = d.conn.CreateTable(input)
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}
}