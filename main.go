package main

import (
	"encoding/json"
	"fmt"
	"os"

	//go get -u github.com/aws/aws-sdk-go
	ses_lib "ses-poc/ses"
	sqs_lib "ses-poc/sqs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
)

type QueueBody struct {
	Email  string `json:"email"`
	S3Key  string `json:"s3_key"`
	FileId string `json:"file_id"`
}

func main() {
	godotenv.Load()

	region := os.Getenv("AWS_REGION")
	queue_url := os.Getenv("SQS_QUEUE_URL")
	bucket_name := os.Getenv("BUCKET_NAME")
	sender_email := os.Getenv("SENDER_EMAIL")

	if region == "" || queue_url == "" || bucket_name == "" || sender_email == "" {
		fmt.Println("Missing Environment Variables")
		return
	}

	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	if err != nil {
		fmt.Println("Error creating session:")
		fmt.Println(err)
		return
	}

	// Create an SES session.
	ses_svc := ses.New(sess)
	sqs_svc := sqs.New(sess)
	s3_svc := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)

	chnMessages := make(chan *sqs.Message, 2)
	go sqs_lib.PollMessages(queue_url, sqs_svc, chnMessages)

	for message := range chnMessages {
		go func(message *sqs.Message) {
			var jsonResult QueueBody
			body := sqs_lib.GetMessage(message)
			err := json.Unmarshal([]byte(*body), &jsonResult)

			if err != nil {
				fmt.Println("Error Unmarshalling JSON")
				fmt.Println(err)
				return
			}

			buff := &aws.WriteAtBuffer{}
			_, err = downloader.Download(buff, &s3.GetObjectInput{
				Bucket: aws.String(bucket_name),
				Key:    aws.String(jsonResult.S3Key),
			})

			s3_svc.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String(bucket_name),
				Key:    aws.String(jsonResult.S3Key),
			})

			if err != nil {
				fmt.Println("Error downloading file:")
				fmt.Println(err)
				return
			}

			ses_lib.SendEmail(ses_svc, sender_email, jsonResult.Email, string(buff.Bytes()), jsonResult.FileId)
			sqs_lib.DeleteMessage(queue_url, sqs_svc, message)
		}(message)
	}
}
