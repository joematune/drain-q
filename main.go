package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func main() {
	queueUrl := flag.String("q", "", "The queue url")
	timeout := flag.Int("t", 20, "How long, in seconds, that the message is hidden from others")
	fileName := flag.String("f", "msgs.txt", "The name of the output file")
	flag.Parse()

	if *queueUrl == "" {
		fmt.Println("You must supply the url of a queue (-q QUEUE_URL)")
		return
	}

	if *timeout < 0 {
		*timeout = 0
	}

	if *timeout > 12*60*60 {
		*timeout = 12 * 60 * 60
	}

	// Open a file for writing
	file, err := os.Create(*fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create a JSON encoder by passing a writer
	encoder := json.NewEncoder(file)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sqs.NewFromConfig(cfg)

	input := &sqs.GetQueueAttributesInput{
		QueueUrl: queueUrl,
		AttributeNames: []types.QueueAttributeName{"ApproximateNumberOfMessages"},
	}

	result, err := client.GetQueueAttributes(context.TODO(), input)
	if err != nil {
		fmt.Println("Got an error getting the queue attributes")
		fmt.Println(err)
		return
	}
	
	approximateMessages, err := strconv.Atoi(result.Attributes["ApproximateNumberOfMessages"])
	if err != nil {
		fmt.Println("Got an error converting string to int")
		fmt.Println(err)
		return
	}
	shouldReceive := approximateMessages > 0

	for shouldReceive {
		receiveMessageInput := &sqs.ReceiveMessageInput{
			MessageAttributeNames: []string{
				string(types.QueueAttributeNameAll),
			},
			QueueUrl:            queueUrl,
			MaxNumberOfMessages: 10,
			VisibilityTimeout:   int32(*timeout),
			WaitTimeSeconds:     20,
		}

		msgResult, err := client.ReceiveMessage(context.TODO(), receiveMessageInput)
		if err != nil {
			fmt.Println("Got an error receiving messages:")
			fmt.Println(err)
			return
		}
		
		if msgResult.Messages == nil {
			fmt.Println("No messages found")
			break
		}

		deleteEntries := []types.DeleteMessageBatchRequestEntry{}

		for _, msg := range msgResult.Messages {
			// Encode and write the data to the file
			err = encoder.Encode(msg)
			if err != nil {
				fmt.Println("Error encoding msg JSON:", err)
				return
			}

			deleteEntries = append(deleteEntries, types.DeleteMessageBatchRequestEntry{
				Id: msg.MessageId,
				ReceiptHandle: msg.ReceiptHandle,
			})
		}

		deleteMessageBatchInput := &sqs.DeleteMessageBatchInput{
			QueueUrl: queueUrl,
			Entries: deleteEntries,
		}

		_, err = client.DeleteMessageBatch(context.TODO(), deleteMessageBatchInput)
		if err != nil {
			fmt.Println("Got an error deleting the messages:")
			fmt.Println(err)
			return
		}
	}
}

