package main
/*
#define CGO_START 0
#define CGO_END 1
*/
import "C"
import (
	"fmt"
	"github.com/Riften/goMesaTracer/tracer"
	"time"
)

var FlagMap = make(map[C.int]string)

func init() {
	fmt.Println("Init flag map.")
	FlagMap[C.CGO_START] = "Trace_Start"
	FlagMap[C.CGO_END] = "Trace_End"
}


// Initialize the routine used for record trace
func init() {
	fmt.Println("Initialize main routine")
	go tracer.GlobalTracer.Start()
}

//export cgoAddTrace
func cgoAddTrace(cgoType C.int) {
	tracer.GlobalTracer.AddRecord(int(cgoType))
}

//export cgoStopAndWait
func cgoStopAndWait() {
	// TODO: End the writer routine and write back logs
	time.Sleep(2 * time.Second) // wait for 2 minute
	return
}

func main() {

	return
}
