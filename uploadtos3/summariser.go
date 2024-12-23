package uploadtos3

import (
	"awsapp/errors"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	comprehendTypes "github.com/aws/aws-sdk-go-v2/service/comprehend/types"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func readContent(body io.Reader) string {
	bytes, err := io.ReadAll(body)
	errors.HandleError(err, "Reading content of S3 bucket record")
	return string(bytes)
}

// func SimulateEvent(bucket, key string) {
// 	event := events.S3Event{
// 		Records: []events.S3EventRecord{
// 			{
// 				EventName: "ObjectCreated:Put",
// 				S3: events.S3Entity{
// 					Bucket: events.S3Bucket{Name: bucket},
// 					Object: events.S3Object{Key: key},
// 				},
// 			},
// 		},
// 	}

// 	HandleRequest(context.Background(), event)

// }

// /*
// Function that is triggered once an upload is made to S3 Bucket.
// */
// func HandleRequest(ctx context.Context, event events.S3Event) {
// 	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
// 	errors.HandleError(err, "HandleRequest")

// 	s3Client := s3.NewFromConfig(cfg)
// 	dynamoClient := dynamodb.NewFromConfig(cfg)
// 	comprehendClient := comprehend.NewFromConfig(cfg)

// 	for idx, record := range event.Records {
// 		fmt.Printf("Processing file %d\n", idx)

// 		bucket := record.S3.Bucket.Name
// 		key := record.S3.Object.Key

// 		resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
// 			Bucket: &bucket,
// 			Key:    &key,
// 		})

// 		if err != nil {
// 			log.Fatalf("Failed to get object information from S3 Bucket. Error: %v", err)
// 		}

// 		defer resp.Body.Close()

// 		text := readContent(resp.Body)
// 		summary := summarize(text, comprehendClient)
// 		storeSummary(key, summary, dynamoClient)
// 	}
// }

func Summarize(content string) string {

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	errors.HandleError(err, "HandleRequest")

	client := comprehend.NewFromConfig(cfg)

	resp, err := client.BatchDetectKeyPhrases(context.TODO(), &comprehend.BatchDetectKeyPhrasesInput{
		TextList:     []string{content},
		LanguageCode: comprehendTypes.LanguageCode(*aws.String("en")),
	})

	if err != nil {
		log.Printf("Failed to summarize content: %v", err)
		return "Summary not available"
	}

	var summary []string
	for _, result := range resp.ResultList {
		for _, keyPhrase := range result.KeyPhrases {
			summary = append(summary, *keyPhrase.Text)
		}
	}
	fmt.Println(strings.Join(summary, ", "))
	return strings.Join(summary, ", ")

}

func storeSummary(key, summary string) {

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	errors.HandleError(err, "HandleRequest")

	dynamoClient := dynamodb.NewFromConfig(cfg)

	_, err = dynamoClient.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String("Papers"),
		Item: map[string]dynamodbTypes.AttributeValue{
			"Key":      &dynamodbTypes.AttributeValueMemberS{Value: key},
			"Summary":  &dynamodbTypes.AttributeValueMemberS{Value: summary},
			"Metadata": &dynamodbTypes.AttributeValueMemberS{Value: "additional metadata here"},
		},
	})

	errors.HandleError(err, "Inserting into DynamoDB")
	fmt.Println("Stored Summary succesfully into DynamoDB")
}
