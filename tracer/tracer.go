package tracer

import "C"

type Tracer struct {
	recv chan string
}

var GlobalTracer Tracer = Tracer{recv: make(chan string, 10)}

func (t Tracer) TryAddTrace(str string) {
	t.recv <- str
}

func (t Tracer) OutCh() chan string {
	return t.recv
}

//export cgoTryAddTrace
func cgoTryAddTrace(str *C.char) {
	gstr := C.GoString(str)
	GlobalTracer.TryAddTrace(gstr)
}