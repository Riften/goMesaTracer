package main
/*
#define CGO_START 0
#define CGO_END 1
#define GLM2_DRAW_BEGIN 2
#define GLM2_DRAW_END 3
#define GLM2_STEP_BEGIN 4
#define GLM2_STEP_END 5
#define GLM2_UPDATE_BEGIN 6
#define GLM2_UPDATE_END 7
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

	tracer.GlobalTracer.W, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println("Error when open output file: ", err.Error())
		panic(err)
	}

	fmt.Println("Init flag map.")
	FlagMap[C.CGO_START] = "Trace_Start"
	FlagMap[C.CGO_END] = "Trace_End"
	FlagMap[C.GLM2_DRAW_BEGIN] = "GLM2_DRAW_BEGIN"
	FlagMap[C.GLM2_DRAW_END] = "GLM2_DRAW_END"
	FlagMap[C.GLM2_STEP_BEGIN] = "GLM2_STEP_BEGIN"
	FlagMap[C.GLM2_STEP_END] = "GLM2_STEP_END"
	FlagMap[C.GLM2_UPDATE_BEGIN] = "GLM2_UPDATE_BEGIN"
	FlagMap[C.GLM2_UPDATE_END] = "GLM2_UPDATE_END"
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
	tracer.GlobalTracer.End()
	time.Sleep(2 * time.Second) // wait for 2 minute
	return
}

func main() {

	return
}
