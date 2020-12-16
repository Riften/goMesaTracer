package tracer

import "C"
import (
	"fmt"
	"io"
	"time"
)

var counter int = 0

type Record struct {
	Counter int // A counter to count the number of lines.
				// It should be set when Record is generated so that we can detect record losing in output file.
	CgoType int	// Type of record.
				// The value of it is defined in main.go as preamble C code.
	TimeStamp int64  // Time when that record happens.
	// [deprecated]
	// SrcDesp and OtherDesp are set as global map.
	// SrcDesp string  // Description of source code.
					// It is the position where record from in most cases.
	// OtherDesp string // Other description. It should be empty in most cases.

	// TODO: Whether we should add duration in record?
}

type Tracer struct {
	Recv chan *Record
	W io.Writer
	Endch chan interface{}
}

// Make the buffer large enough so that the tracer would not block the main thread.
var GlobalTracer Tracer = Tracer{
	Recv: make(chan *Record, 1000),
	Endch: make(chan interface{}),
}

func (t Tracer) AddRecord(cgoType int) {
	t.Recv <- &Record{
		Counter:   counter,
		CgoType:   cgoType,
		TimeStamp: time.Now().UnixNano(),
	}
	counter += 1
}

func (t Tracer) WriteRaw(r *Record) {
	// counter cgotype timestamp
	_, err := t.W.Write([]byte(fmt.Sprintf("%d %d %d\n", r.Counter, r.CgoType, r.TimeStamp)))
	if err != nil {
		fmt.Println("Error when write record back: ", err.Error())
	}
}

// Call Start in a separate goroutine.
// Note:
//		The priority of case is higher than default in go.
//		So Recv and Endch works as two priority queue.
func (t Tracer) Start() {
	fmt.Println("Tracer routine start.")
	for {
		select {
		case data := <- t.Recv:
			t.WriteRaw(data)
		default:
			select {
			case data := <- t.Recv:
				t.WriteRaw(data)
			case <- t.Endch:
				fmt.Println("Tracer routine end.")
				return
			}
		}
	}
}

// Call End when there is nothing more to trace.
func (t Tracer) End() {
	t.Endch <- struct {}{}
}