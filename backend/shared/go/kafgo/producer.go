package kafgo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

// How to use this. Create a kafka producer. User Send() to send payloads. The payload should be a struct with json tags, cause it will get marshaled.

type KafkaProducer struct {
	client *kgo.Client
}

// seeds are used for finding the server, just as many kafka ip's you have
func NewKafkaProducer(seeds []string) (producer *KafkaProducer, close func(), err error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, nil, err
	}
	kfc := &KafkaProducer{
		client: cl,
	}
	return kfc, cl.Close, nil
}

// TODO batch sends instead of doing one by one
// Send sends payload(s) to the specified topic
func (kfc *KafkaProducer) Send(ctx context.Context, topic string, payload ...any) error {
	records := make([]*kgo.Record, len(payload))
	for i, p := range payload {
		bytes, err := json.Marshal(p)
		if err != nil {
			return err
		}
		records[i] = &kgo.Record{Topic: topic, Value: bytes}
	}
	results := kfc.client.ProduceSync(ctx, records...)
	if results.FirstErr() != nil {
		return fmt.Errorf("failed to produce %w", results.FirstErr())
	}
	return nil
}
