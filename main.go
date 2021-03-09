package main

import "C"
import (
	"fmt"
	"github.com/Riften/goMesaTracer/common"
	"github.com/Riften/goMesaTracer/tracer"
	"os"
	"time"
)

func init() {
	var err error

	mesaCmdOnly := os.Getenv("MESA_TRACE_CMD_ONLY")
	mesaNoOut := os.Getenv("MESA_TRACE_NO_OUT")
	if mesaCmdOnly != "" {
		tracer.GlobalTracer.OutCmdOnly = true
		fmt.Println("Cgo trace would be print out to command line only")
	} else if mesaNoOut != "" {
		tracer.GlobalTracer.NoOut = true
		fmt.Println("Cgo trace would have no output.")
	} else {
		tracer.GlobalTracer.OutCmdOnly = false
		filePath := os.Getenv("MESA_TRACE_OUT")
		if filePath == "" {
			fmt.Println("No out file specified, use default out file ", common.DefaultOutFile)
			fmt.Println("You can set the out file by os env MESA_TRACE_OUT")
			filePath =common.DefaultOutFile
		}

		tracer.GlobalTracer.W, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println("Error when open output file: ", err.Error())
			panic(err)
		}
	}

	tracer.GlobalTracer.FetchFlagName = func(cgoType int) string {
		return common.FlagMap[cgoType]
	}
	fmt.Println("Initialize main routine")
	go tracer.GlobalTracer.Start()
}

//export cgoAddTrace
func cgoAddTrace(counter C.int, funcName *C.char) {
	tracer.GlobalTracer.AddRecord(int(counter), C.GoString(funcName));
}

//export cgoStopAndWait
func cgoStopAndWait() {
	// TODO: End the writer routine and write back logs
	tracer.GlobalTracer.End()
	time.Sleep(2 * time.Second) // wait for 2 minute
	return
}

func main() {
	err := Run()
	if err != nil {
		fmt.Println("Error in cmd: ", err)
	}
	return
}
