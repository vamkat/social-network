package kafgo

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"social-network/shared/go/batching"
	"social-network/shared/go/ct"
	tele "social-network/shared/go/telemetry"
	"sort"
	"strconv"
	"sync"
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

// func TestKafka() error {
// 	err := TestKafkaProducer()
// 	if err != nil {
// 		return err
// 	}
// 	return TestKafkaConsumer()
// }

func HugeKafkaTest() {

	randomTopic := "test_" + fmt.Sprint(rand.IntN(10000))

	ctx := context.Background()
	messages := []string{}
	for i := range 200 {
		messages = append(messages, fmt.Sprint(i))
	}

	wg := sync.WaitGroup{}

	wg.Go(func() {
		err := TestKafkaProducer(messages, randomTopic)
		if err != nil {
			tele.Error(context.Background(), "huge kafka test producer @1", "error", err.Error())
			return
		}
	})

	err := TestKafkaConsumer(messages, randomTopic)
	if err != nil {
		tele.Error(ctx, "failed consumer test @1", "error", err.Error())
		return
	}
	tele.Info(context.Background(), "success!")

}

func TestKafkaProducer(messages []string, topic string) error {
	ctx := context.Background()

	producer, _, err := NewKafkaProducer([]string{"kafka:9092"})
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	handler := func(messages [][]byte) error {
		tele.Info(ctx, "batcher sending, @1 @2", "from", string(messages[0]), "to", string(messages[len(messages)-1]))
		err := producer.Send(ctx, topic, messages...)
		if err != nil {
			return fmt.Errorf("failed to send messages, err: %w", err)
		}
		return nil
	}

	batchInput, errChan := batching.Batcher(ctx, handler, time.Millisecond*100, 1000)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		for _, msg := range messages {

			select {
			case batchInput <- []byte(msg):
			case <-ctx.Done():
				return
			}
			time.Sleep(time.Millisecond * 10)
		}
		ctx.Done()
	})

	wg.Go(func() {
		for err := range errChan {
			if err != nil {
				tele.Error(ctx, "found error in producer huge test @1", "error", err.Error())
				ctx.Done()
				return
			}
		}
	})

	wg.Wait()
	return nil
}

func TestKafkaConsumer(messages []string, topic string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer, err := NewKafkaConsumer([]string{"kafka:9092"}, "test")
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	consumer.WithCommitBuffer(50)

	ch, err := consumer.RegisterTopic(ct.KafkaTopic(topic))
	if err != nil {
		return fmt.Errorf("failed to register topic: %w", err)
	}

	found := []string{}
	mu := sync.Mutex{}

	consumerRoutine := func() {
		me := rand.IntN(10000)
		for {
			tele.Info(ctx, "@1 start loop", "me", me)
			select {
			case record := <-ch:

				msg := string(record.rec.Value)
				tele.Info(ctx, "@1 [CONSUMER] received @2", "me", me, "msg", msg)

				dur := min(
					rand.IntN(500),
					rand.IntN(500),
					rand.IntN(500),
					rand.IntN(500),
					rand.IntN(500),
					rand.IntN(500),
					rand.IntN(500),
					rand.IntN(500),
					rand.IntN(500),
					rand.IntN(500),
					rand.IntN(500),
					rand.IntN(500),
					rand.IntN(500),
				)
				time.Sleep(time.Millisecond * time.Duration(dur))

				tele.Info(ctx, "@1 [CONSUMER] committing @2 after waiting for @3ms", "me", me, "msg", msg, "time", dur)
				err := record.Commit(ctx)
				if err != nil {
					tele.Error(context.Background(), "consumer huge test @1", "error", err.Error())
					return
				}
				tele.Info(ctx, "@1 [CONSUMER] committed @2", "me", me, "msg", msg)

				mu.Lock()
				tele.Info(ctx, "@1 [CONSUMER] adding @2", "me", me, "msg", msg)
				found = append(found, msg)

				if len(found) == len(messages) {
					mu.Unlock()
					return
				}
				mu.Unlock()
			case <-ctx.Done():
				tele.Info(ctx, "@1 [CONSUMER] DONE", "me", me)
				return
			}

		}
	}

	wg := sync.WaitGroup{}
	for range 3 {
		wg.Go(func() { consumerRoutine() })
	}

	closeAll, err := consumer.StartConsuming(ctx)
	if err != nil {
		tele.Error(ctx, "failed to start consuming @1", "error", err.Error())
	}
	defer closeAll()

	wg.Wait()

	sort.Slice(found, func(i, j int) bool {
		jVal, _ := strconv.Atoi(found[j])
		iVal, _ := strconv.Atoi(found[i])
		return iVal < jVal
	})

	for i, msg := range found {
		tele.Info(ctx, "comparing -> @1 and @2", "found", msg, "expected", fmt.Sprint(i))
		if msg != fmt.Sprint(i) {
			return errors.New(fmt.Sprint("final expected matching test failed! at index: ", i, " f:", msg, " is not same as:", fmt.Sprint(i)))
		}
	}

	if len(messages) != len(found) {
		return fmt.Errorf("incorrect expected found count msgs: %d, found: %d", len(messages), len(found))
	}
	tele.Info(ctx, "Testing success! Waiting before closing.")
	time.Sleep(time.Second * 5)
	tele.Info(ctx, "CONSUMER ENDED")
	return nil
}
