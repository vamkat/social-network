package kafgo

import (
	"context"
	"errors"
	"os"
	"social-network/shared/go/ct"
	tele "social-network/shared/go/telemetry"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type committerData struct {
	TopicChannel       chan *Record
	CommitChannel      chan *Record
	StartOffset        chan int64
	CommitRoutineReady bool
	ExpectedOffsets    []int64
}

func newCommitRoutine() *committerData {
	return &committerData{
		TopicChannel:  make(chan *Record),
		CommitChannel: make(chan *Record),
		StartOffset:   make(chan int64),
	}
}

// To explain the idea of this package. You create a kafka consumer. You then call RegisterTopic() to register topics which gives you a channel to listen to.
// You will receive a Record{} from this channel. You get your data from the Data() method, process it, and once you're done you Commit(), and that's all you gotta do!

type kafkaConsumer struct {
	topics            []string
	seeds             []string
	group             string
	context           context.Context
	client            *kgo.Client
	committerDataMap  map[string]*committerData
	allOffsetsSorted  bool
	offtsetsOkCount   int
	commitBuffer      int
	topicBuffer       int
	cancel            func()
	isConsuming       bool
	weAreShuttingDown bool
}

//TODO check if consumer and producer work gracefully after kafka restarts

// seeds are used for finding the server, just as many kafka ip's you have
// enter all the topics you want to consume
// enter you group identifier
func NewKafkaConsumer(seeds []string, group string) (*kafkaConsumer, error) {
	if len(seeds) == 0 || group == "" {
		return nil, errors.New("NewKafkaConsumer: bad arguments")
	}

	kfc := &kafkaConsumer{
		seeds:            seeds,
		group:            group,
		committerDataMap: make(map[string]*committerData),
		commitBuffer:     1000,
		topicBuffer:      5000,
	}

	return kfc, nil
}

func (kfc *kafkaConsumer) WithCommitBuffer(size int) *kafkaConsumer {
	if kfc.isConsuming {
		panic("don't mess with the consumer while it's consuming!")
	}
	kfc.commitBuffer = size
	return kfc
}

func (kfc *kafkaConsumer) WithTopicBuffer(size int) *kafkaConsumer {
	kfc.topicBuffer = size
	return kfc
}

var ErrFetch = errors.New("error when fetching")
var ErrConsumerFunc = errors.New("consumer function error")

// RegisterTopic registers a topic for consumption and returns a channel to receive records
// Recommended buffer size is above 100, should probably be in the thousands, set to 0 for default
func (kfc *kafkaConsumer) RegisterTopic(topic ct.KafkaTopic) (<-chan *Record, error) {
	if kfc.isConsuming {
		panic("you can't register topics in the middle of consuming!")
	}

	_, ok := kfc.committerDataMap[string(topic)]
	if ok {
		panic("you've passed duplicate topics")
	}

	kfc.topics = append(kfc.topics, string(topic))
	topicChannel := make(chan *Record, kfc.topicBuffer)
	kfc.committerDataMap[string(topic)] = newCommitRoutine()

	return topicChannel, nil
}

// StartConsuming sets some stuff up and begin the consumption routines
func (kfc *kafkaConsumer) StartConsuming(ctx context.Context) (func(), error) {
	var err error
	//making the actual client, cause it needs to be created after all topics have been registered
	kfc.client, err = kgo.NewClient(
		kgo.SeedBrokers(kfc.seeds...),
		kgo.ConsumerGroup(kfc.group),
		kgo.ConsumeTopics(kfc.topics...),
		kgo.DisableAutoCommit(),
		kgo.AllowAutoTopicCreation(),
		kgo.Balancers(kgo.CooperativeStickyBalancer()),
	)
	if err != nil {
		return nil, err
	}

	ctx, cancelContext := context.WithCancel(ctx)
	closeAll := func() {
		kfc.client.Close()
		cancelContext()
	}

	kfc.context = ctx
	kfc.cancel = closeAll

	if err := kfc.validateBeforeStart(); err != nil {
		return nil, err
	}

	kfc.actuallyStartConsuming()

	return closeAll, nil
}

// actuallyStartConsuming actually does the consumption
func (kfc *kafkaConsumer) actuallyStartConsuming() {
	kfc.isConsuming = true

	// commitChannels are listened to by a routine for that topic's record commits
	// after the handlers are done processing the record they call commit, and these routines get informed
	for _, topic := range kfc.topics {
		commitData := kfc.committerDataMap[topic]
		commitData.CommitChannel = make(chan *Record, kfc.commitBuffer)
		commitData.StartOffset = make(chan int64)
		go kfc.commitRoutine(commitData)
	}

	go func() {
		timer := time.NewTimer(time.Second)
		defer timer.Stop()
		for {
			select {
			case <-kfc.context.Done():
				kfc.shutdownProcedure(false)
				return
			default:
				fetches := kfc.client.PollFetches(kfc.context)
				if errs := fetches.Errors(); len(errs) > 0 {
					// All errors are retried internally when fetching, but non-retriable errors are
					// returned from polls so that users can notice and take action.
					tele.Info(context.Background(), "fetch error: @1", "error", errs)
					kfc.shutdownProcedure(true)
					return
				}

				// We can iterate through a record iterator...
				iter := fetches.RecordIter()
				for !iter.Done() {
					record := iter.Next()

					committerData := kfc.committerDataMap[record.Topic]
					committerData.ExpectedOffsets = append(committerData.ExpectedOffsets, record.Offset)
					//since the commit routines need to commit the records in the right order
					//we need to find the smallest record offset and let them know
					if !committerData.CommitRoutineReady {
						committerData.StartOffset <- record.Offset
						committerData.CommitRoutineReady = true
						kfc.offtsetsOkCount++
						if kfc.offtsetsOkCount == len(kfc.topics) {
							kfc.allOffsetsSorted = true
						}
					}

					Record, err := newRecord(record, committerData.CommitChannel)
					if err != nil {
						//think what to do
						tele.Info(context.Background(), "failed to create record")
						continue
					}

					timer.Reset(time.Second * 5)
					select {
					case <-timer.C:
						tele.Info(context.Background(), "SLOW CHANNEL DETECTED")
						kfc.shutdownProcedure(true)
						tele.Info(context.Background(), "SLOW CHANNEL error: ")
						return
					case committerData.TopicChannel <- Record:
					}
				}
			}
		}
	}()
}

// shutdownProcedure handles the shutdown process of the consumer
// makes sure everything is closed properly
// is indempodent so it can be called again without a problem
func (kfc *kafkaConsumer) shutdownProcedure(thereIsSomethingWrong bool) {

	//to ensure idempotensy
	if kfc.weAreShuttingDown {
		return
	}

	kfc.weAreShuttingDown = true

	if thereIsSomethingWrong {
		tele.Error(kfc.context, "SHUTTING DOWN KAFKA CONSUMER")
	} else {
		tele.Info(kfc.context, "Shutting down kafka consumer")
	}

	//cancelling the context, both of the kafka inner client, and this packages context
	kfc.cancel()

	//closing all topic channels, so that no more record are sent to handlers
	for _, committerData := range kfc.committerDataMap {
		close(committerData.TopicChannel)
	}

	//ranging over the topics again to drain them and discard the records
	for _, committerData := range kfc.committerDataMap {
		for range committerData.TopicChannel {
		}
	}

	// timer := time.NewTimer(time.Second * 10) //TODO use me!

	//committing any remaining commits
	//TODO do this!

	if thereIsSomethingWrong {
		os.Exit(1)
	}
}

//TODO make separate commit channels and routine for each topic
//TODO handle out of order commits...
//TODO batch commits, small batches
//TODO add detection trap whenc committing offsets out of order

// commitRoutine listens to the commitChannel and commits records as they come in
func (kfc *kafkaConsumer) commitRoutine(data *committerData) {
	defer tele.Info(context.Background(), "COMMIT WATCHER ROUTINE CLOSING DEFERRED")
	//we wait for the consumer loop to give us the offset of the first record if receives
	nextOffset := <-data.StartOffset
	tooEarlyOffsets := make(map[int64]struct{}, 30)
	for {
		select {
		case <-kfc.context.Done():
			tele.Info(context.Background(), "COMMIT WATCHER ROUTINE CLOSING DUE TO CONTEXT")
			return

		case record := <-data.CommitChannel:
			tooEarlyOffsets[record.rec.Offset] = struct{}{}

			_, nextExists := tooEarlyOffsets[data.ExpectedOffsets[0]]
			for nextExists {
				//TODO pool records here instead of doing them one by one
				ctx, cancel := context.WithTimeout(kfc.context, time.Second*2) //TODO is this the correct context?
				defer cancel()
				err := kfc.client.CommitRecords(ctx, record.rec) //TODO is this the correct context?

				if err != nil {
					tele.Info(context.Background(), "COMMIT ERROR FOUND @1", "error", err.Error()) //TODO think what needs to be done here
					kfc.shutdownProcedure(true)                                                    //TODO this is excessive, but not sure what else to do? other than retry?
				}

				data.ExpectedOffsets = data.ExpectedOffsets[1:]
				if len(data.ExpectedOffsets) == 0 {
					break
				}
				_, nextExists = tooEarlyOffsets[nextOffset]
			}

		}
	}
}

func (kfc *kafkaConsumer) validateBeforeStart() error {
	if kfc.context == nil {
		return errors.New("nil context")
	}
	select {
	case <-kfc.context.Done():
		return errors.New("context already canceled")
	default:
	}

	if kfc.isConsuming {
		return errors.New("consumer already started")
	}
	if kfc.weAreShuttingDown {
		return errors.New("consumer is shutting down")
	}

	if len(kfc.seeds) == 0 {
		return errors.New("no seeds configured")
	}
	if kfc.group == "" {
		return errors.New("no consumer group configured")
	}

	if len(kfc.topics) == 0 {
		return errors.New("no topics registered")
	}
	if len(kfc.committerDataMap) != len(kfc.topics) {
		return errors.New("topic/channel mismatch")
	}

	seen := make(map[string]struct{}, len(kfc.topics))
	for _, t := range kfc.topics {
		if _, ok := seen[t]; ok {
			return errors.New("duplicate topic")
		}
		seen[t] = struct{}{}
		ch, ok := kfc.committerDataMap[t]
		if !ok || ch == nil {
			return errors.New("missing topic channel")
		}
	}

	if kfc.commitBuffer <= 0 {
		return errors.New("invalid commit buffer")
	}

	return nil
}
