package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
)

type rawTrace struct {
	count int
	cgoType int
	nano int64
}

type stackTrace struct {
	raw *rawTrace
	depth int	// depth maintain the call stack depth
	//order int	// order maintain the in-stack order
				// [deprecated]: this can be done by raw.counter
	duration int64
	startTime int64
	endTime int64
}

type stackStatistic struct {
	calls []callStatistic
	startTime int64
	endTime int64
	busyTime int64
}

type callStatistic struct {
	count int
	totalDuration int64
	avgDuration int64
}

func newStackStatistic() *stackStatistic {
	_calls := make([]callStatistic, totalFlag)
	for _, c := range _calls {
		c.totalDuration = 0
		c.count = 0
		c.avgDuration = 0
	}
	return &stackStatistic{
		calls:     _calls,
		startTime: 0,
		endTime:   0,
		busyTime:  0,
	}
}

func scanTrace(r io.Reader, handler func(trace *rawTrace)) {
	scanner := bufio.NewScanner(r)
	var n_scan int
	var err error
	for scanner.Scan() {
		tmpRawTrace := &rawTrace{}
		n_scan, err = fmt.Sscanf(scanner.Text(), "%d %d %d", &tmpRawTrace.count, &tmpRawTrace.cgoType, &tmpRawTrace.nano)
		if err != nil {
			fmt.Println("Error when scan line: ", err)
		} else if n_scan==0 {
			fmt.Println("Error: Sscanf parsed no param.")
		}

		handler(tmpRawTrace)
	}
}

type stacker struct {
	dest io.Writer
	src io.Reader
	stack *lls.Stack

	//writeBuf []*stackTrace
	//bufSize int
	writeBuf *stackWriteBuf
	stStatistic *stackStatistic
}

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
/*
type flushWriteBuf func(buf *stackWriteBuf)

func (st *stacker) flushBuf(buf *stackWriteBuf) {

}
 */

func (st *stacker) writeLn(str string) {
	_, err := st.dest.Write([]byte(str+"\n"))
	if err != nil {
		fmt.Println("Error when write to stacker dest: ", err)
	}
}

func (st *stacker) writeTrace(tr *stackTrace) {
	_, err := st.dest.Write([]byte(strings.Repeat("\t", tr.depth)))
	if err != nil {
		fmt.Println("Error when write to stacker dest: ", err)
	}
	_, err = st.dest.Write([]byte(fmt.Sprintf("%s %d\n", strings.TrimSuffix(FlagMap[tr.raw.cgoType], "_BEGIN"), tr.duration)))
	if err != nil {
		fmt.Println("Error when write to stacker dest: ", err)
	}
}

func (st *stacker) analyze() {
	scanTrace(st.src, st.handleRawTrace)
}

func (st *stacker) handleRawTrace(r *rawTrace) {
	if r.cgoType == 0 {
		st.writeLn(FlagMap[0])
		return
	}
	if r.cgoType == 1 {
		st.writeLn(FlagMap[1])
		return
	}

	// monitor stack process:
	if r.count % 100 == 0 {
		fmt.Printf("Process line %d\r", r.count)
	}

	if r.cgoType % 2 == 0 {
		peekt := st.peekTrace()
		var dep int
		if peekt==nil {
			dep = 0
		} else {
			dep = peekt.depth + 1
		}

		st.stack.Push(&stackTrace{
			raw:      r,
			depth:    dep,
			duration: 0, 	// duration is computed when pop stack
			startTime: r.nano,
		})
	} else {
		for	{
			peekt := st.peekTrace()
			if peekt == nil {
				// When the stack is empty
				fmt.Printf("Stack empty when get %d %s %d\n", r.count, FlagMap[r.cgoType], r.nano)
				break
			} else if peekt.raw.cgoType != r.cgoType-1 {
				// When the top trace mismatch
				st.stack.Pop()
				fmt.Printf("Stack mismatch\n\tget %d %s %d\n\tpeek %d %s %d\n",
					r.count, FlagMap[r.cgoType], r.nano, peekt.raw.count, FlagMap[peekt.raw.cgoType], peekt.raw.nano)
				continue
			} else {
				// Trace matched
				st.stack.Pop()
				peekt.duration = r.nano - peekt.raw.nano
				peekt.endTime = r.nano
				st.writeBuf.add(peekt)
				//st.writeBuf[st.bufSize] = peekt
				//st.bufSize += 1

				if st.stack.Empty() {
					st.writeBuf.flush(st.writeTrace)
				}
				break
			}
		}
	}
}

func (st *stacker) peekTrace() *stackTrace {
	if st.stack.Empty() {
		return nil
	}
	t, _ := st.stack.Peek()
	return t.(*stackTrace)
}

// Used to analyze the call stack from trace
func cmdStack(inputPath string, outPath string) error {
	inFile, err := os.OpenFile(inputPath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Error when open input trace file: ", err)
		return err
	}

	if outPath == "" {
		inputSuffix := path.Ext(inputPath)
		withoutSuffix := strings.TrimSuffix(inputPath, inputSuffix)
		outPath = withoutSuffix + "_stack.txt"
		fmt.Println("No outPath specified. Out to default file ", outPath)
	}
	fmt.Println("Call stack will be writen to ", outPath)

	outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error when open output stack file: ", err)
		return err
	}
	// do rw
	st := &stacker{
		dest:  outFile,
		src:   inFile,
		stack: lls.New(),
		writeBuf: &stackWriteBuf{
			buf:  make([]*stackTrace, 1000),
			size: 0,
		},
	}

	st.analyze()

	inFile.Close()
	outFile.Close()
	return nil
}

