package enum

/*
 enum TraceFlag {Trace_Start, Trace_End}
 */
import "C"
import "fmt"

var FlagMap = make(map[C.TraceFlag]string)

//export cgoGetFlagString
func cgoPrintFlagString(flag C.TraceFlag) {
	str, ok := FlagMap[flag]
	if !ok {
		fmt.Println("No such flag")
	} else {
		fmt.Println(str)
	}
}

func init() {
	fmt.Println("Init flag map.")
	FlagMap[C.Trace_Start] = "Trace_Start"
	FlagMap[C.Trace_End] = "Trace_End"
}
