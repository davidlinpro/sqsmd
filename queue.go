package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type QueueStats struct {
	Name          string `json:"name"`     // Name of queue
	Messages      int    `json:"messages"` // ApproximateNumberOfMessages - Returns the approximate number of visible messages in a queue.
	MessagesMax   int    `json:"messages_max"`
	Delayed       int    `json:"delayed"` // ApproximateNumberOfMessagesDelayed - Returns the approximate number of messages that are waiting to be added to the queue.
	DelayedMax    int    `json:"delayed_max"`
	NotVisible    int    `json:"not_visible"` // ApproximateNumberOfMessagesNotVisible - Returns the approximate number of messages that have not timed-out and aren't deleted.
	NotVisibleMax int    `json:"not_visible_max"`
}

func (qs *QueueStats) SetMax() {
	if qs.Messages > qs.MessagesMax {
		qs.MessagesMax = qs.Messages
	}
	if qs.Delayed > qs.DelayedMax {
		qs.DelayedMax = qs.Delayed
	}
	if qs.NotVisible > qs.NotVisibleMax {
		qs.NotVisibleMax = qs.NotVisible
	}
}

func (qs *QueueStats) Copy(q QueueStats) {
	qs.Messages = q.Messages
	qs.Delayed = q.Delayed
	qs.NotVisible = q.NotVisible
	qs.SetMax()
}

func processQueueStats() {
	// connect to SQS
	creds := credentials.NewStaticCredentials(Config.AWSACCESSKEY, Config.AWSSECRETKEY, "")
	cfg := aws.NewConfig().WithRegion(Config.AWSREGION).WithCredentials(creds)
	sess := session.Must(session.NewSessionWithOptions(session.Options{Config: *cfg}))
	svc := sqs.New(sess)

	// get list of queues
	result, err := svc.ListQueues(nil)
	if err != nil {
		log.Fatal(err)
	}

	// save only custom queues with prefix
	var customQueues []string
	for _, url := range result.QueueUrls {
		if url == nil {
			log.Println("invalid nil queue")
			continue
		}
		u := strings.Replace(*url, Config.SQSBASEURL, "", 1)
		if Config.SQSPREFIX != "" && !strings.HasPrefix(u, Config.SQSPREFIX) {
			continue
		}
		customQueues = append(customQueues, u)
	}

	// get initial stats for queues and save to cache
	GetStatsAndSave(svc, customQueues)

	// then update stats every 10 seconds (should be a variable)
	for range time.Tick(10 * time.Second) {
		go GetStatsAndSave(svc, customQueues)
	}
}

func GetStatsAndSave(svc *sqs.SQS, customQueues []string) {
	mainAttributes := []*string{
		aws.String("ApproximateNumberOfMessages"),
		aws.String("ApproximateNumberOfMessagesDelayed"),
		aws.String("ApproximateNumberOfMessagesNotVisible"),
	}

	for _, name := range customQueues {
		q := QueueStats{Name: name}
		params := &sqs.GetQueueAttributesInput{
			QueueUrl:       aws.String(Config.SQSBASEURL + name),
			AttributeNames: mainAttributes,
		}
		resp, _ := svc.GetQueueAttributes(params)
		q.Messages, _ = strconv.Atoi(*resp.Attributes["ApproximateNumberOfMessages"])
		q.Delayed, _ = strconv.Atoi(*resp.Attributes["ApproximateNumberOfMessagesDelayed"])
		q.NotVisible, _ = strconv.Atoi(*resp.Attributes["ApproximateNumberOfMessagesNotVisible"])
		CacheSet(q)
		LogIt(q)
	}
}
