package stack

import "sort"

// NOT THREAD SAFE!!
type stackWriteBuf struct {
	buf []*stackTrace
	size int
	// out io.Writer
}

func (wf *stackWriteBuf) add(st *stackTrace) {
	wf.buf[wf.size] = st
	wf.size += 1
}

// Let stackWriteBuf sortable
func (wf *stackWriteBuf) Len() int {
	return wf.size
}

func (wf *stackWriteBuf) Swap(i, j int) {
	wf.buf[i], wf.buf[j] = wf.buf[j], wf.buf[i]
}

func (wf *stackWriteBuf) Less(i, j int) bool {
	return wf.buf[i].raw.count < wf.buf[j].raw.count
}

func (wf *stackWriteBuf) flush(handleWrite func(trace *stackTrace)) {
	sort.Sort(wf)
	for i:=0 ; i<wf.size; i+=1{
		//fmt.Println("...", FlagMap[curTrace.raw.cgoType])
		handleWrite(wf.buf[i])
	}
	wf.size = 0
}

// Used as stack buffer when running in `stack full` mode
// stackFull is a re
type stackFullBuf struct {

}
