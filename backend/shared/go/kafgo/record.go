package kafgo

import (
	"context"
	"errors"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Record is a type that helps with commiting after processing a record
// AFTER PROCESSING THE RECORD MAKE SURE TO COMMIT!!!
type Record struct {
	rec           *kgo.Record
	commitChannel chan<- (*kgo.Record)
}

var ErrBadArgs = errors.New("bro, you passed bad arguments")

// newRecord creates a new Record instance
func newRecord(record *kgo.Record, commitChannel chan<- (*kgo.Record)) (*Record, error) {
	if record == nil {
		return nil, fmt.Errorf("%w record: %v", ErrBadArgs, record)
	}
	return &Record{
		rec:           record,
		commitChannel: commitChannel,
	}, nil
}

// Data returns the data payload
func (r *Record) Data() []byte {
	if r.rec == nil {
		//log?
		return []byte("no data found")
	}
	return r.rec.Value
}

var ErrEmptyRecord = errors.New("empty record")

// TODO commit must return a confirmation

// Commit marks the record as processed in the Kafka client.
// MAKE SURE THIS IS AT THE END OF A TRANSACTION, DONT BE COMMITING THINGS YOU LATER UNDO!!
func (r *Record) Commit(ctx context.Context) error {
	if r.rec == nil {
		return ErrEmptyRecord
	}
	select {
	case r.commitChannel <- r.rec:
	case <-ctx.Done():
		// optionally log or ignore
	}

	return nil
}
