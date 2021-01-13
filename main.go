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
#define GLX_SWAP_BUFFERS_BEGIN 30
#define GLX_SWAP_BUFFERS_END 31
#define GLX_DRI3_FLUSH_DRIAWABLE_BEGIN 32
#define GLX_DRI3_FLUSH_DRIAWABLE_END 33
#define ST_FLUSH_BEGIN 34
#define ST_FLUSH_END 35
#define GAL_TC_FLUSH_BEGIN 36
#define GAL_TC_FLUSH_END 37
#define GAL_DRI_FLUSH_BEGIN 38
#define GAL_DRI_FLUSH_END 39
#define GAL_DRI_THRO_PRE_BEGIN 40
#define GAL_DRI_THRO_PRE_END 41
#define MESA_CLEAR_BEGIN 42
#define MESA_CLEAR_END 43
#define MESA_CREATE_VAO_BEGIN 44
#define MESA_CREATE_VAO_END 45
#define MESA_GEN_VAO_BEGIN 46
#define MESA_GEN_VAO_END 47
#define MESA_BIND_VAO_BEGIN 48
#define MESA_BIND_VAO_END 49
#define MESA_GEN_VBO_BEGIN 50
#define MESA_GEN_VBO_END 51
#define MESA_CREATE_VBO_BEGIN 52
#define MESA_CREATE_VBO_END 53
#define MESA_BIND_VBO_BEGIN 54
#define MESA_BIND_VBO_END 55
#define MESA_VBO_DATA_BEGIN 56
#define MESA_VBO_DATA_END 57
#define MESA_VERTEX_ATTRIB_POINTER_BEGIN 58
#define MESA_VERTEX_ATTRIB_POINTER_END 59
#define GLM2_GLX_SWAP_BEGIN 60
#define GLM2_GLX_SWAP_END 61
#define GLM2_CANVAS_UPDATE_BEGIN 62
#define GLM2_CANVAS_UPDATE_END 63
#define ZINK_FENCE_FINISH_BEGIN 64
#define ZINK_FENCE_FINISH_END 65
#define VK_WAIT_FENCES_BEGIN 66
#define VK_WAIT_FENCES_END 67
#define ZINK_START_BATCH_BEGIN 68
#define ZINK_START_BATCH_END 69
#define ZINK_END_BATCH_BEGIN 70
#define ZINK_END_BATCH_END 71
#define VK_QUEUE_SUBMIT_BEGIN 72
#define VK_QUEUE_SUBMIT_END 73
#define VK_END_COMMAND_BUFFER_BEGIN 74
#define VK_END_COMMAND_BUFFER_END 75
#define VK_BEGIN_COMMAND_BUFFER_BEGIN 76
#define VK_BEGIN_COMMAND_BUFFER_END 77
#define ZINK_CREATE_FRAMEBUF_BEGIN 78
#define ZINK_CREATE_FRAMEBUF_END 79
#define ZINK_SHADER_COMPILE_BEGIN 80
#define ZINK_SHADER_COMPILE_END 81
#define TEXIMAGE_BEGIN 82
#define TEXIMAGE_END 83

#define MESA_BUFFER_DATA 101
#define MESA_TEXIMAGE 102
*/
import "C"
import (
	"fmt"
	"github.com/Riften/goMesaTracer/common"
	"github.com/Riften/goMesaTracer/tracer"
	"os"
	"time"
)



//var FlagMap = make([]string, common.TotalFlag)

func init() {
	var err error

	mesaCmdOnly := os.Getenv("MESA_TRACE_CMD_ONLY")
	if mesaCmdOnly != "" {
		tracer.GlobalTracer.OutCmdOnly = true
		fmt.Println("Cgo trace would be print out to command line only")
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

//export cgoAddDetail
func cgoAddDetail(cgoType C.int, detail C.longlong) {
	tracer.GlobalTracer.AddDetail(int(cgoType), int64(detail))
}

//export cgoAddDoubleInt
func cgoAddDoubleInt(cgoType C.int, i1 C.int, i2 C.int) {
	tracer.GlobalTracer.AddDetail(int(cgoType), int64(i1 * i2))
}

func main() {
	err := Run()
	if err != nil {
		fmt.Println("Error in cmd: ", err)
	}
	return
}
