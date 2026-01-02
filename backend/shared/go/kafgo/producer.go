package kafgo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

type kafkaProducer struct {
	client *kgo.Client
}

// seeds are used for finding the server, just as many kafka ip's you have
func NewKafkaProducer(ctx context.Context, seeds []string) (*kafkaProducer, func(), error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
	)
	if err != nil {
		return nil, nil, err
	}
	kfc := &kafkaProducer{
		client: cl,
	}
	return kfc, cl.Close, nil
}

var ErrProduceFail = errors.New("failed to produce")

func (kfc *kafkaProducer) Send(ctx context.Context, topic string, payload ...any) error {
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
