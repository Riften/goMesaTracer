package main
/*
#define CGO_START 0
#define CGO_END 1
*/
import "C"
import (
	"fmt"
	"github.com/Riften/goMesaTracer/tracer"
	"os"
	"time"
)

const defaultOutFile = "mesa_trace_raw.csv"

var FlagMap = make(map[C.int]string)

func init() {
	var err error

	filePath := os.Getenv("MESA_TRACE_OUT")
	if filePath == "" {
		fmt.Println("No out file specified, use default out file ", defaultOutFile)
		fmt.Println("You can set the out file by os env MESA_TRACE_OUT")
		filePath = defaultOutFile
	}

	tracer.GlobalTracer.W, err = os.OpenFile(os.Getenv("MESA_TRACE_OUT"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 644)
	if err != nil {
		fmt.Println("Error when open output file: ", err.Error())
		panic(err)
	}

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
