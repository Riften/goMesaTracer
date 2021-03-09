package tracer

import "C"
import (
	"fmt"
	"io"
	"time"
)

// Record ...
type Record struct {
	Counter   int
	FuncName  string
	TimeStamp int64
}

// Tracer ...
type Tracer struct {
	Recv       chan *Record
	W          io.Writer
	Endch      chan interface{}
	OutCmdOnly bool // Whether print the trace to command line only.
	// It would be true if the os environment MESA_TRACE_CMD_ONLY is not empty.
	NoOut bool // Run with no output.
	// It would be true if the os environment MESA_TRACE_NO_OUT is not empty.

	FetchFlagName func(int) string // Used to fetch name string of cgoType.
}

// GlobalTracer ...
// Make the buffer large enough so that the tracer would not block the main thread.
var GlobalTracer Tracer = Tracer{
	Recv:  make(chan *Record, 1000),
	Endch: make(chan interface{}),
}

// AddRecord ...
func (t Tracer) AddRecord(counter int, funcName string) {
	t.Recv <- &Record{
		Counter:   counter,
		FuncName:  funcName,
		TimeStamp: time.Now().UnixNano(),
	}
}

// WriteRaw ...
func (t Tracer) WriteRaw(r *Record) {
	var err error
	_, err = t.W.Write([]byte(fmt.Sprintf("%d %s %d\n", r.Counter, r.FuncName, r.TimeStamp)))
	if err != nil {
		fmt.Println("Error when write record back: ", err.Error())
	}
}

// WriteCmd ...
func (t Tracer) WriteCmd(r *Record) {
	fmt.Printf("%d %s %d\n", r.Counter, r.FuncName, r.TimeStamp)
}

// DoNothing ...
func (t Tracer) DoNothing(r *Record) {

}

// Start ...
// Call Start in a separate goroutine.
// Note:
//		The priority of case is higher than default in go.
//		So Recv and Endch works as two priority queue.
func (t Tracer) Start() {
	fmt.Println("Tracer routine start.")
	var writeFunc func(*Record)
	if t.OutCmdOnly {
		writeFunc = t.WriteCmd
	} else if t.NoOut {
		writeFunc = t.DoNothing
	} else {
		writeFunc = t.WriteRaw
	}
	for {
		select {
		case data := <-t.Recv:
			writeFunc(data)
		default:
			select {
			case data := <-t.Recv:

				// TODO: Remove this when it seems that 1000 is enough for channel buffer
				if data.Counter%100 == 0 {
					if len(t.Recv) > 500 {
						fmt.Println("Warning: half of tracer buffer is full.")
					}
				}
				writeFunc(data)

			case <-t.Endch:
				fmt.Println("Tracer routine end.")
				return
			}
		}
	}
}

// End ...
// Call End when there is nothing more to trace.
func (t Tracer) End() {
	t.Endch <- struct{}{}
}
