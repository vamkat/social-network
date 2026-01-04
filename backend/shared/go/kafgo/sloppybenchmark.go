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
