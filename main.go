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
#define ZINK_DRAW_BEGIN 8
#define ZINK_DRAW_END 9
#define MESA_SET_DRAW_VAO_BEGIN 10
#define MESA_SET_DRAW_VAO_END 11
#define FLUSH_FOR_DRAW_BEGIN 12
#define FLUSH_FOR_DRAW_END 13
#define MESA_DRAW_ARRAYS_BEGIN 14
#define MESA_DRAW_ARRAYS_END 15
#define ST_DRAW_VBO_BEGIN 16
#define ST_DRAW_VBO_END 17
#define ST_PREPARE_DRAW_BEGIN 18
#define ST_PREPARE_DRAW_END 19
#define CSO_DRAW_VBO_BEGIN 20
#define CSO_DRAW_VBO_END 21
#define GLM2_BUILD_RENDER_VBO_BEGIN 22
#define GLM2_BUILD_RENDER_VBO_END 23
#define GLM2_BUILD_RENDER_ARRAY_BEGIN 24
#define GLM2_BUILD_RENDER_ARRAY_END 25
#define GLM2_BUILD_LOAD_PROJECTION_BEGIN 26
#define GLM2_BUILD_LOAD_PROJECTION_END 27
#define GLM2_BUILD_LOAD_NORMAL_BEGIN 28
#define GLM2_BUILD_LOAD_NORMAL_END 29
*/
import "C"
import (
	"fmt"
	"github.com/Riften/goMesaTracer/tracer"
	"os"
	"time"
)

const defaultOutFile = "mesa_trace_raw.csv"

var FlagMap = make([]string, 100)

func init() {
	var err error

	filePath := os.Getenv("MESA_TRACE_OUT")
	if filePath == "" {
		fmt.Println("No out file specified, use default out file ", defaultOutFile)
		fmt.Println("You can set the out file by os env MESA_TRACE_OUT")
		filePath = defaultOutFile
	}

	tracer.GlobalTracer.W, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
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
	FlagMap[C.ZINK_DRAW_BEGIN] = "ZINK_DRAW_BEGIN"
	FlagMap[C.ZINK_DRAW_END] = "ZINK_DRAW_END"
	FlagMap[C.MESA_SET_DRAW_VAO_BEGIN] = "MESA_SET_DRAW_VAO_BEGIN"
	FlagMap[C.MESA_SET_DRAW_VAO_END] = "MESA_SET_DRAW_VAO_END"
	FlagMap[C.FLUSH_FOR_DRAW_BEGIN] = "FLUSH_FOR_DRAW_BEGIN"
	FlagMap[C.FLUSH_FOR_DRAW_END] = "FLUSH_FOR_DRAW_END"
	FlagMap[C.MESA_DRAW_ARRAYS_BEGIN] = "MESA_DRAW_ARRAYS_BEGIN"
	FlagMap[C.MESA_DRAW_ARRAYS_END] = "MESA_DRAW_ARRAYS_END"
	FlagMap[C.ST_DRAW_VBO_BEGIN] = "ST_DRAW_VBO_BEGIN"
	FlagMap[C.ST_DRAW_VBO_END] = "ST_DRAW_VBO_END"
	FlagMap[C.ST_PREPARE_DRAW_BEGIN] = "ST_PREPARE_DRAW_BEGIN"
	FlagMap[C.ST_PREPARE_DRAW_END] = "ST_PREPARE_DRAW_END"
	FlagMap[C.GLM2_LOAD_PROJECTION_BEGIN] = "GLM2_LOAD_PROJECTION_BEGIN"
	FlagMap[C.GLM2_LOAD_PROJECTION_END] = "GLM2_LOAD_PROJECTION_END"
	FlagMap[C.GLM2_BUILD_RENDER_VBO_BEGIN] = "GLM2_BUILD_RENDER_VBO_BEGIN"
	FlagMap[C.GLM2_BUILD_RENDER_VBO_END] = "GLM2_BUILD_RENDER_VBO_END"
	FlagMap[C.GLM2_BUILD_RENDER_ARRAY_BEGIN] = "GLM2_BUILD_RENDER_ARRAY_BEGIN"
	FlagMap[C.GLM2_BUILD_RENDER_ARRAY_END] = "GLM2_BUILD_RENDER_ARRAY_END"
	FlagMap[C.GLM2_BUILD_LOAD_PROJECTION_BEGIN] = "GLM2_BUILD_LOAD_PROJECTION_BEGIN"
	FlagMap[C.GLM2_BUILD_LOAD_PROJECTION_END] = "GLM2_BUILD_LOAD_PROJECTION_END"
	FlagMap[C.GLM2_BUILD_LOAD_NORMAL_BEGIN] = "GLM2_BUILD_LOAD_NORMAL_BEGIN"
	FlagMap[C.GLM2_BUILD_LOAD_NORMAL_END] = "GLM2_BUILD_LOAD_NORMAL_END"
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
	err := Run()
	if err != nil {
		fmt.Println("Error in cmd: ", err)
	}
	return
}
