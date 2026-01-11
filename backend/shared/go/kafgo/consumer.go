package kafgo

import (
	"context"
	"errors"
	"os"
	"social-network/shared/go/ct"
	tele "social-network/shared/go/telemetry"
	"sync"
	"sync/atomic"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type topicData struct {
	TopicChannel  chan *Record //the consumer uses this to send messages to whoever needs messages
	CommitChannel chan *Record //used by Record.Commit() so that the message handler can communicate to the commit routine
	ExpectedIds   []uint64     //these are id's the consumer generated that correspond to the records it fetches, it's needed so that committing records happens in the same order as they arrived
	Mutex         *sync.Mutex
}

func newCommitRoutine() *topicData {
	return &topicData{
		TopicChannel:  make(chan *Record),
		CommitChannel: make(chan *Record),
		Mutex:         &sync.Mutex{},
	}
}

// To explain the idea of this package. You create a kafka consumer. You then call RegisterTopic() to register topics which gives you a channel to listen to.
// You will receive a Record{} from this channel. You get your data from the Data() method, process it, and once you're done you Commit(), and that's all you gotta do!

type kafkaConsumer struct {
	topics            []string
	seeds             []string //use by the franz package to connect to the cluster, put as many as you have
	group             string   //consumer group means that messages as distributed to the members of the same group
	context           context.Context
	client            *kgo.Client           //what actually does operation to kafka
	topicDataMap      map[string]*topicData //data needed for each topic
	commitBuffer      int                   //how big the commit channel buffer
	topicBuffer       int                   //how big the topic channels' buffer
	cancel            func()
	isConsuming       bool
	weAreShuttingDown bool
}

//ALL TODOs
//set up proper limits on how many records this consumer can have at the same time.
//test more that shutdown is graceful enough
//test that records don't get lost when service gets restarted
//add comitting remaining documents on shutdown
//add batching to commits
//make a trap for detecting out of order documents
//add major retrying behavior instead of defaulting to a shutdown

// Seeds are used for finding the server, just as many kafka ip's you have.
//
// Enter all the topics you want to consume.
//
// Enter your group identifier because within the same group messages get spread evenly (at least it tries to).
//
// Usage:
//
//			//create consumer
//			consumer, err := kafgo.NewKafkaConsumer([]string{"localhost:9092"}, "chat")
//			if err != nil {
//	     	tele.Error(ctx, "wtf")
//			}
//
//			//register topics. You will be given a
//			//channel to listen to for messages of that topic
//			memberChannel, err := consumer.RegisterTopic("notifications")
//			alertChannel, err = consumer.RegisterTopic("alerts")
//
//			//then activate the consumption
//			close, err := consumer.StartConsuming(ctx)
//			if err != nil {
//				tele.Error(ctx, "wtf")
//			}
//			defer close()
func NewKafkaConsumer(seeds []string, group string) (*kafkaConsumer, error) {
	if len(seeds) == 0 || group == "" {
		return nil, errors.New("NewKafkaConsumer: bad arguments")
	}

	kfc := &kafkaConsumer{
		seeds:        seeds,
		group:        group,
		topicDataMap: make(map[string]*topicData),
		commitBuffer: 1000,
		topicBuffer:  5000,
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
	if kfc.isConsuming {
		panic("don't mess with the consumer while it's consuming!")
	}
	kfc.topicBuffer = size
	return kfc
}

var ErrFetch = errors.New("error when fetching")
var ErrConsumerFunc = errors.New("consumer function error")

// RegisterTopic registers a topic for consumption and returns a channel to receive records
func (kfc *kafkaConsumer) RegisterTopic(topic ct.KafkaTopic) (<-chan *Record, error) {
	if kfc.isConsuming {
		panic("you can't register topics in the middle of consuming!")
	}

	_, ok := kfc.topicDataMap[string(topic)]
	if ok {
		panic("you've passed duplicate topics")
	}

	kfc.topics = append(kfc.topics, string(topic))
	committerData := newCommitRoutine()
	kfc.topicDataMap[string(topic)] = committerData

	return committerData.TopicChannel, nil
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
		commitData := kfc.topicDataMap[topic]
		commitData.CommitChannel = make(chan *Record, kfc.commitBuffer)
		go kfc.commitRoutine(commitData)
	}

	go func() {
		// This id will be used for identifying orders so that they can
		// be committed in the same order as they arrived.
		// Kafka's offset can't be used cause it gets, seemingly, randomly reset to 0...
		var monotonicIds atomic.Uint64

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
					for i, err := range errs {
						tele.Error(kfc.context, "fetch @1 error: @2", "number", i, "error", err.Err.Error())
					}

					kfc.shutdownProcedure(true)
					return
				}
				tele.Info(kfc.context, "consumer fetch successful")

				// We can iterate through a record iterator...
				iter := fetches.RecordIter()
				for !iter.Done() {
					newId := monotonicIds.Add(1)
					record := iter.Next()
					committerData := kfc.topicDataMap[record.Topic]

					//since the commit routines need to commit the records in the right order
					//we need keep track of what records have arrived, this is then read by commit routine
					tele.Info(kfc.context, "consumer before mutex")
					committerData.Mutex.Lock()
					committerData.ExpectedIds = append(committerData.ExpectedIds, newId)
					committerData.Mutex.Unlock()
					tele.Info(kfc.context, "consumer after mutex")

					Record, err := newRecord(kfc.context, record, committerData.CommitChannel, newId)
					if err != nil {
						//think what to do
						tele.Error(context.Background(), "failed to create record @1", "error", err.Error())
						continue
					}

					timer.Reset(time.Second * 5)
					tele.Info(kfc.context, "consumer before timer select")
					select {
					case <-timer.C:
						tele.Error(context.Background(), "SLOW CHANNEL DETECTED")
						kfc.shutdownProcedure(true)
						return
					case committerData.TopicChannel <- Record:
						tele.Info(kfc.context, "consumer give record to topic channel")
					}
					tele.Info(kfc.context, "consumer after timer select")
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
	for _, committerData := range kfc.topicDataMap {
		close(committerData.TopicChannel)
	}

	//ranging over the topics again to drain them and discard the records
	for _, committerData := range kfc.topicDataMap {
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

//TODO batch commits, small batches
//TODO add detection trap whenc committing monotinicIds out of order

// commitRoutine listens to the commitChannel and commits records as they come in
// it batches records to reduce spamming too many requests to kafka <- TODO
// since records need to be committed in order, it collects record monotinicIds,
// and by using the expected monotinicIds provided by the consumer routine we only commit the next expected record
// HOPEFULLY THE RECORD HANDLERS DO THEIR BEST EFFORT TO PROCESS IT, and if it's corrupted then they should commit it to remove it from circulation
// if a handler is incapable of of processing the record, then it needs to close the channel
// so that the shutdown operation can begin and restart the pod
func (kfc *kafkaConsumer) commitRoutine(data *topicData) {
	foundRecords := make(map[uint64]*Record, 30)
	for {
		select {
		case <-kfc.context.Done():
			tele.Info(context.Background(), "COMMIT WATCHER ROUTINE CLOSING DUE TO CONTEXT .Done()")
			return

		case newRecord, ok := <-data.CommitChannel:
			if !ok {
				tele.Error(kfc.context, "commitroutine: shutting down due to bad channel")
				//if not ok, means channel was closed by record handler, which means there's a critical problem and pod needs to be restarted
				kfc.shutdownProcedure(true)
				return
			}

			//add new record to collection of monotonicIds
			foundRecords[newRecord.monotinicId] = newRecord
			tele.Info(kfc.context, "new record found of @1. current @2, and expected monoIds @3", "monoId", newRecord.monotinicId, "count", len(foundRecords), "monoIdsLen", len(data.ExpectedIds))

			//check if the next expected monoId (assigned by consumer routine) if available
			data.Mutex.Lock()
			nextMonoId := data.ExpectedIds[0]
			data.Mutex.Unlock()

			record, nextExists := foundRecords[nextMonoId]
			combo := 0
			if !nextExists {
				tele.Info(kfc.context, "next record of @1 not found, skipping...", "monoId", nextMonoId)
			}
			for nextExists {
				tele.Info(kfc.context, "inside nextExists")
				combo++
				//TODO pool records here instead of doing them one by one
				ctx, cancel := context.WithTimeout(kfc.context, time.Second*5) //TODO is this the correct context?
				err := kfc.client.CommitRecords(ctx, record.rec)               //TODO is this the correct context?
				if err != nil {
					tele.Error(context.Background(), "COMMIT ERROR FOUND @1", "error", err.Error()) //TODO think what needs to be done here
					kfc.shutdownProcedure(true)                                                     //TODO this is excessive, but not sure what else to do? other than retry?
					cancel()
					return
				}
				cancel()

				//the handler of the record is waiting for confirmation that we received it before committing the transaction
				//so lets confirm it
				record.confirmChannel <- struct{}{}

				//clean map from committed record

				delete(foundRecords, nextMonoId)

				//time to check if the next found expected monoId is available too
				data.Mutex.Lock()
				data.ExpectedIds = data.ExpectedIds[1:]
				if len(data.ExpectedIds) == 0 {
					data.Mutex.Unlock()
					if combo > 1 {
						tele.Info(kfc.context, "COMBO: @1", "count", combo)
					}
					break
				}
				nextMonoId = data.ExpectedIds[0]
				record, nextExists = foundRecords[nextMonoId]
				data.Mutex.Unlock()
			}
			if combo > 1 {
				tele.Info(kfc.context, "COMBO: @1", "count", combo)
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
	if len(kfc.topicDataMap) != len(kfc.topics) {
		return errors.New("topic/channel mismatch")
	}

	seen := make(map[string]struct{}, len(kfc.topics))
	for _, t := range kfc.topics {
		if _, ok := seen[t]; ok {
			return errors.New("duplicate topic")
		}
		seen[t] = struct{}{}
		ch, ok := kfc.topicDataMap[t]
		if !ok || ch == nil {
			return errors.New("missing topic channel")
		}
	}

	if kfc.commitBuffer <= 0 {
		return errors.New("invalid commit buffer")
	}

	return nil
}
