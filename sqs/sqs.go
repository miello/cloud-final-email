package sqs_lib

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func PollMessages(queue_url string, sqsSvc *sqs.SQS, chn chan<- *sqs.Message) {

	for {
		output, err := sqsSvc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queue_url),
			MaxNumberOfMessages: aws.Int64(2),
			WaitTimeSeconds:     aws.Int64(15),
		})

		if err != nil {
			fmt.Println("failed to fetch sqs message %v", err)
		}

		for _, message := range output.Messages {
			chn <- message
		}

	}
}

func GetMessage(msg *sqs.Message) *string {
	fmt.Println("RECEIVING MESSAGE >>> ")

	return msg.Body
}

func DeleteMessage(queue_url string, sqsSvc *sqs.SQS, msg *sqs.Message) {
	sqsSvc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queue_url),
		ReceiptHandle: msg.ReceiptHandle,
	})
}
