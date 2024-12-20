package rapidbug

import (
	"fmt"
	"io"
)

type FakeReader struct {
	events     []ReadEvent
	eventIndex int
	dataIndex  int
}

type ReadEvent struct {
	Data string
	Err  error
}

func (re ReadEvent) String() string {
	return fmt.Sprintf("ReadEvent{Data: %q, Err: %v}", re.Data, re.Err)
}

// NewFakeReader creates a new FakeReader that can perform
// short reads and occassionally return 0, nil even if the
// input buffer to Read is not empty.
//
// Pre-condition: events must not contain io.EOF.
func NewFakeReader(events []ReadEvent) FakeReader {
	return FakeReader{
		events:     events,
		eventIndex: 0,
		dataIndex:  0,
	}
}

var _ io.Reader = &FakeReader{}

func (r *FakeReader) Read(p []byte) (_ int, err error) {
	var eof error
	if r.eventIndex == len(r.events) {
		eof = io.EOF
	}
	if len(p) == 0 {
		return 0, eof
	}
	if eof != nil {
		return 0, eof
	}
	event := r.events[r.eventIndex]
	ncopy := min(len(event.Data[r.dataIndex:]), len(p))
	lo, hi := r.dataIndex, r.dataIndex+ncopy
	copy(p, event.Data[lo:hi])
	if hi == len(event.Data) {
		r.dataIndex = 0
		r.eventIndex++
	} else {
		r.dataIndex += ncopy
	}
	return ncopy, event.Err
}
