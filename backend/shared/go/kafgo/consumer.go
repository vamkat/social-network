package kafgo

import (
	"context"
	"errors"
	"fmt"
	"social-network/shared/go/ct"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type kafkaConsumer struct {
	context       context.Context
	client        *kgo.Client
	topicChannels map[string]chan<- *Record
}

// seeds are used for finding the server, just as many kafka ip's you have
// enter the topics you want to consume, if any
// enter you group identifier
func NewKafkaConsumer(ctx context.Context, seeds []string, topics []string, group string) (*kafkaConsumer, func(), error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topics...),
		kgo.DisableAutoCommit(),
	)

	if err != nil {
		return nil, nil, err
	}

	kfc := &kafkaConsumer{
		context:       ctx,
		client:        cl,
		topicChannels: make(map[string]chan<- *Record, 3),
	}

	return kfc, cl.Close, nil
}

var ErrFetch = errors.New("error when fetching")
var ErrConsumerFunc = errors.New("consumer function error")

func (kfc *kafkaConsumer) RegisterTopic(ctx context.Context, topic ct.KafkaTopic) (<-chan *Record, error) {
	outputChan := make(chan *Record)
	kfc.topicChannels[string(topic)] = outputChan
	return outputChan, nil
}

func (kfc *kafkaConsumer) StartConsuming(ctx context.Context) func() {
	commitChannel := make(chan (*kgo.Record))
	newCtx, cancel := context.WithCancel(ctx)
	go func() {
		for {
			select {
			case <-newCtx.Done():
				return
			default:
				fetches := kfc.client.PollFetches(ctx)
				if errs := fetches.Errors(); len(errs) > 0 {
					// All errors are retried internally when fetching, but non-retriable errors are
					// returned from polls so that users can notice and take action.
					kfc.shutdownProcedure()
					return
				}

				// We can iterate through a record iterator...
				iter := fetches.RecordIter()
				for !iter.Done() {
					record := iter.Next()

					Record, err := newRecord(record, commitChannel)
					if err != nil {
						//thing what to do
						fmt.Println("failed to create record")
						continue
					}

					//TODO set a timer for detecting choking channels
					timer := time.NewTimer(time.Second)
					select {
					case <-timer.C:
						fmt.Print("SLOW CHANNEL DETECTED")
					case kfc.topicChannels[record.Topic] <- Record:
					}
				}
			}
		}
	}()

	return cancel
}

func (kfc *kafkaConsumer) shutdownProcedure() {
	//stop consuming loop
	//close all channels immediately
	//drain commit channel
}

func (kfc *kafkaConsumer) commitRoutine(commitChannel <-chan Record) {
	for {
		select {
		case <-kfc.context.Done():
			return
		case record := <-commitChannel:
			fmt.Println("record:", record)
			//TODO pool records here instead of doing them one by one
			err := kfc.client.CommitRecords(context.Background(), record.rec) //TODO pick correct context
			if err != nil {
				fmt.Println("ERRO FOUND") //TODO think what needs to be done here
			}

		}
	}
}
