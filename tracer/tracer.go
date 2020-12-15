package tracer

import "C"

type Tracer struct {
	recv chan string
}

var GlobalTracer Tracer = Tracer{recv: make(chan string, 10)}

func (t Tracer) tryAddTrace(str string) {
	t.recv <- str
}

//export cgoTryAddTrace
func cgoTryAddTrace(str *C.char) {
	gstr := C.GoString(str)
	tracer
}