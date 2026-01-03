package kafgo

import (
	"context"
	"fmt"
	tele "social-network/shared/go/telemetry"
	"time"
)

func Spam() {

	ctx := context.Background()
	tele.Info(ctx, "START OF SPAM")
	//
	//

	// consumer, err := NewKafkaConsumer([]string{"kafka:9092"}, "test")
	// if err != nil {
	// 	tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	// }
	// consumer.WithCommitBuffer(10)
	// consumer.WithCommitBuffer(20)

	// ch, err := consumer.RegisterTopic("test_topic")
	// if err != nil {
	// 	tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	// }

	// go func() {
	// 	defer func() {
	// 		if rec := recover(); rec != nil {
	// 			tele.Error(ctx, "panic occured in consumer loop!")
	// 		}
	// 	}()
	// 	for {
	// 		tele.Info(ctx, "consume loop start")
	// 		for record := range ch {
	// 			time.Sleep(time.Millisecond * 5)
	// 			record.Commit(ctx)
	// 		}
	// 		time.Sleep(time.Second * 3)
	// 	}

	// }()

	// _, err = consumer.StartConsuming(ctx)
	// if err != nil {
	// 	tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	// }
	// tele.Info(ctx, "started consuming")

	//
	//
	//
	//
	producer, _, err := NewKafkaProducer([]string{"kafka:9092"})
	if err != nil {
		tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
	}

	type X struct {
		Name string `json:"name"`
	}

	fmt.Println(producer)

	go func() {
		for i := range 1000000 {
			dur := 15
			time.Sleep(time.Millisecond * time.Duration(dur))
			tele.Info(ctx, "sending after waiting for: @1", "millis", dur)
			err := producer.Send(ctx, "test_topic", X{fmt.Sprint("alex:", i, " and slept for: ", dur, "ms sadfl")})
			if err != nil {
				tele.Error(ctx, "KAFKA ERROR: @1", "error", err.Error())
			}

		}
	}()

	tele.Info(ctx, "END OF SPAM")
	time.Sleep(time.Minute)
}
