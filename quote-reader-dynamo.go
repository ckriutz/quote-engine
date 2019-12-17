package main

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"fmt"
)

const awsRegion = "us-east-2"
const dynamoTable = "quotes"

func checkDynamoForQuotesTable() bool {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// create the input configuration instance
	input := &dynamodb.ListTablesInput{}

	// Get the list of tables
	result, err := svc.ListTables(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return false
	}

	for _, n := range result.TableNames {
		if *n == dynamoTable {
			fmt.Println("Quotes table found in AWS. Looking good.")
			return true
		}
	}

	fmt.Println("Quotes table not found, this is concerning.")
	return false
}

func getCountOfQuotes() int {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(dynamoTable),
		Key: map[string]*dynamodb.AttributeValue{
			"QuoteId": {
				S: aws.String("count"),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	cnt, err := strconv.Atoi(*result.Item["count"].S)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	fmt.Println("Database indicates there are " + *result.Item["count"].S + " quotes in the database.")
	return cnt
}

func getQuoteFromDynamo(id int) DynamoQuote {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(dynamoTable),
		Key: map[string]*dynamodb.AttributeValue{
			"QuoteId": {
				S: aws.String(strconv.Itoa(id)),
			},
		},
	})

	if err != nil {
		fmt.Println(err.Error())
		return DynamoQuote{}
	}

	quote := DynamoQuote{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &quote)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return quote
}

func addQuoteToDynamo(q DynamoQuote) bool {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(q)
	if err != nil {
		fmt.Println("Got error marshalling new DynamoQuote:")
		fmt.Println(err.Error())
		return false
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dynamoTable),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return false
	}

	updateQuoteCount()
	fmt.Println("Successfully added quote id '" + q.QuoteId + "' to table " + dynamoTable)
	return true
}

func updateQuote(q DynamoQuote) bool {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#text": aws.String("text"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":author": {
				S: aws.String(q.Author),
			},
			":t": {
				S: aws.String(q.Text),
			},
		},
		TableName: aws.String(dynamoTable),
		Key: map[string]*dynamodb.AttributeValue{
			"QuoteId": {
				S: aws.String(q.QuoteId),
			},
		},
		UpdateExpression: aws.String("SET author = :author, #text = :t"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println("Successfully updated quote id '" + q.QuoteId + "' to table " + dynamoTable)
	return true
}

func updateQuoteCount() bool {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion)},
	)

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	currentCount := getCountOfQuotes() + 1

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeNames: map[string]*string{
			"#count": aws.String("count"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":count": {
				S: aws.String(strconv.Itoa(currentCount)),
			},
		},
		TableName: aws.String(dynamoTable),
		Key: map[string]*dynamodb.AttributeValue{
			"QuoteId": {
				S: aws.String("count"),
			},
		},
		UpdateExpression: aws.String("SET #count = :count"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println("Successfully incremented quote count.")
	return true
}
