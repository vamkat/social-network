package kafgo

import (
	"context"
	"errors"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Record is a type that helps with commiting after processing a record
type Record struct {
	rec           *kgo.Record
	commitChannel chan<- (*kgo.Record)
}

var ErrBadArgs = errors.New("bad arguments passed")

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
		return []byte{}
	}
	return r.rec.Value
}

// Commit marks the record as processed in the Kafka client.
// MAKE SURE THIS IS PART OF A TRANSACTION, DONT BE COMMITING THINGS YOU LATER UNDO!!
func (r *Record) Commit(ctx context.Context) {
	if r.rec == nil {
		return
	}

	select {
	case r.commitChannel <- r.rec:
	case <-ctx.Done():
		// optionally log or ignore
	}
}
