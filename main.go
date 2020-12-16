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

//export cgoAddString
func cgoAddString(str *C.char) {
	gstr := C.GoString(str)
	tracer.GlobalTracer.TryAddTrace(gstr)
}

//export cgoTestEnum
func cgoTestEnum(flagEnum C.int) {
	str, ok := FlagMap[flagEnum]
	if !ok {
		fmt.Println("Unknown flag enum ", flagEnum)
	} else {
		tracer.GlobalTracer.TryAddTrace(str + " " + "testTrace")
	}
}

// Initialize the routine used for record trace
func init() {
	fmt.Println("Initialize main routine")
	go func() {
		defer func() {
			fmt.Println("Routine End")
		}()
		fmt.Println("Routine start")
		var recved string
		var ok bool
		outCh := tracer.GlobalTracer.OutCh()
		for {
			recved, ok = <-outCh
			if !ok {
				fmt.Println("Tracer channel closed!")
				return
			} else {
				fmt.Println(recved)
			}
		}
	}()
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
