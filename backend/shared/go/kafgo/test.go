package kafgo

import (
	"context"
	"fmt"
	"math/rand/v2"
	tele "social-network/shared/go/telemetry"
	"sync/atomic"
	"time"
)

func SpamConsumer() {

	ctx := context.Background()
	tele.Info(ctx, "START OF CONSUMER SPAM")

	//
	//
	//
	//
	//
	consumer, err := NewKafkaConsumer([]string{"kafka:9092"}, "test")
	if err != nil {
		tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	}

	consumer.WithCommitBuffer(1000)

	ch, err := consumer.RegisterTopic("test_topic")
	if err != nil {
		tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	}

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				tele.Error(ctx, "panic occured in consumer loop!")
			}
		}()
		for {
			a := atomic.Int64{}
			tele.Info(ctx, "consume loop start")
			for record := range ch {
				tele.Info(ctx, "found record! starting go routine")
				go func() {
					dur := 5 + min(rand.IntN(100), rand.IntN(100), rand.IntN(100), rand.IntN(100), rand.IntN(100), rand.IntN(200))
					time.Sleep(time.Millisecond * time.Duration(dur))
					tele.Info(ctx, "handler attempting to commit")
					err := record.Commit(ctx)
					tele.Info(ctx, "get confirmationf or: @1 #@2", "offset", record.rec.Offset, "a", a.Load())
					a.Add(1)
					if err != nil {
						tele.Error(ctx, "error with commit! @1", "error", err.Error())
					}
					tele.Info(ctx, "handler finished commit after waiting for: @1", "millis", dur)
				}()
			}
			time.Sleep(time.Second * 3)
		}

	}()

	_, err = consumer.StartConsuming(ctx)
	if err != nil {
		tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	}
	tele.Info(ctx, "started consuming")
	tele.Info(ctx, "END OF CONSUMER SPAM")
}

type X struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func SpamProducer() {

	ctx := context.Background()
	tele.Info(ctx, "START OF Producer SPAM")
	//
	//
	//
	//
	//
	producer, _, err := NewKafkaProducer([]string{"kafka:9092"})
	if err != nil {
		tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	}

	fmt.Println(producer)
	go func() {
		for i := range 200000 {
			dur := 30
			time.Sleep(time.Millisecond * time.Duration(dur))
			tele.Info(ctx, "sending after waiting for: @1", "millis", dur)
			err := producer.Send(ctx, "test_topic", []byte(fmt.Sprint("alex:", i, " and slept for: ", dur, "ms count: ", i)))
			if err != nil {
				tele.Error(ctx, "KAFKA SEND ERROR: @1", "error", err.Error())
				return
			}

		}
	}()

	tele.Info(ctx, "END OF Producer SPAM")
}

func TestKafka() error {
	err := TestKafkaProducer()
	if err != nil {
		return err
	}
	return TestKafkaConsumer()
}

func TestKafkaProducer() error {
	ctx := context.Background()

	producer, _, err := NewKafkaProducer([]string{"kafka:9092"})
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}

	expected := []string{
		"message_0", "message_1", "message_2", "message_3", "message_4",
		"message_5", "message_6", "message_7", "message_8", "message_9",
		"message_10", "message_11", "message_12", "message_13", "message_14",
		"message_15", "message_16", "message_17", "message_18", "message_19",
	}

	for _, msg := range expected {
		err := producer.Send(ctx, "test_topic", []byte(msg))
		if err != nil {
			return fmt.Errorf("failed to send message %s: %w", msg, err)
		}
		time.Sleep(time.Millisecond * 20)
	}

	return nil
}

func TestKafkaConsumer() error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	consumer, err := NewKafkaConsumer([]string{"kafka:9092"}, "test")
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	consumer.WithCommitBuffer(1000)

	ch, err := consumer.RegisterTopic("test_topic")
	if err != nil {
		return fmt.Errorf("failed to register topic: %w", err)
	}

	expected := []string{
		"message_0", "message_1", "message_2", "message_3", "message_4",
		"message_5", "message_6", "message_7", "message_8", "message_9",
		"message_10", "message_11", "message_12", "message_13", "message_14",
		"message_15", "message_16", "message_17", "message_18", "message_19",
	}

	received := make([]string, 0, 20)
	resultCh := make(chan error, 1)

	go func() {
		for record := range ch {
			tele.Info(ctx, "received record from consumer channel!")

			msg := string(record.rec.Value)
			received = append(received, msg)

			tele.Info(ctx, "commiting!")
			err := record.Commit(ctx)
			if err != nil {
				resultCh <- fmt.Errorf("commit failed for %s: %w", msg, err)
				return
			}

			if len(received) == 20 {
				for i := range expected {
					if received[i] != expected[i] {
						resultCh <- fmt.Errorf("order violation at index %d: expected %s, got %s", i, expected[i], received[i])
						return
					}
				}
				resultCh <- nil
				return
			}
		}
	}()

	_, err = consumer.StartConsuming(ctx)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	select {
	case result := <-resultCh:
		return result
	case <-ctx.Done():
		return fmt.Errorf("timeout waiting for messages, received %d/20", len(received))
	}
}
