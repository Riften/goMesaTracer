package stack

import (
	"bufio"
	"fmt"
	"github.com/Riften/goMesaTracer/common"
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
	"io"
	"os"
	"path"
	"strings"
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
	_calls := make([]callStatistic, common.TotalFlag)
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
		} else {
			handler(tmpRawTrace)
		}
	}
}

type stacker struct {
	dest io.Writer
	src io.Reader
	stack *lls.Stack

	writeBuf *stackWriteBuf
	full bool
	stStatistic *stackStatistic

	firstCallToCompare string
	secondCallToCompare string
}

func newStacker(src io.Reader, dest io.Writer, _full bool) *stacker {
	st := &stacker{
		dest:  dest,
		src:   src,
		full: _full,
		stack: lls.New(),
		writeBuf: &stackWriteBuf{
			buf:  make([]*stackTrace, 1000),
			size: 0,
		},
	}
	return st
}

func (st *stacker) setFirstCompareCall(c string) *stacker {
	st.firstCallToCompare = c
	return st
}

func (st *stacker) setSecondCompareCall(c string) *stacker {
	st.secondCallToCompare = c
	return st
}

func (st *stacker) writeLn(str string) {
	_, err := st.dest.Write([]byte(str+"\n"))
	if err != nil {
		fmt.Println("Error when write to stacker dest: ", err)
	}
}

func (st *stacker) writeTrace(tr *stackTrace) {
	_, err := st.dest.Write([]byte(strings.Repeat("  ", tr.depth)))
	if err != nil {
		fmt.Println("Error when write to stacker dest: ", err)
	}
	_, err = st.dest.Write([]byte(fmt.Sprintf("%s %d\n", strings.TrimSuffix(common.FlagMap[tr.raw.cgoType], "_BEGIN"), tr.duration)))
	if err != nil {
		fmt.Println("Error when write to stacker dest: ", err)
	}
}

func (st *stacker) analyze() {
	if st.full {
		fmt.Println("Stacker is running in full mode, which may be a bit slow but safer.")
	}

	scanTrace(st.src, st.handleRawTrace)
}

// used to handle each raw trace when running in full mode.
func (st *stacker) handleRawTraceFull(r *rawTrace) {
	if r.cgoType == 0 {
		st.writeLn(common.FlagMap[0])
		return
	}
	if r.cgoType == 1 {
		st.writeLn(common.FlagMap[1])
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
				fmt.Printf("Stack empty when get %d %s %d\n", r.count, common.FlagMap[r.cgoType], r.nano)
				break
			} else if peekt.raw.cgoType != r.cgoType-1 {
				// When the top trace mismatch
				st.stack.Pop()
				fmt.Printf("Stack mismatch\n\tget %d %s %d\n\tpeek %d %s %d\n",
					r.count, common.FlagMap[r.cgoType], r.nano, peekt.raw.count, common.FlagMap[peekt.raw.cgoType], peekt.raw.nano)
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

func (st *stacker) handleRawTrace(r *rawTrace) {
	if r.cgoType == 0 {
		st.writeLn(common.FlagMap[0])
		return
	}
	if r.cgoType == 1 {
		st.writeLn(common.FlagMap[1])
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
				fmt.Printf("Stack empty when get %d %s %d\n", r.count, common.FlagMap[r.cgoType], r.nano)
				break
			} else if peekt.raw.cgoType != r.cgoType-1 {
				// When the top trace mismatch
				st.stack.Pop()
				fmt.Printf("Stack mismatch\n\tget %d %s %d\n\tpeek %d %s %d\n",
					r.count, common.FlagMap[r.cgoType], r.nano, peekt.raw.count, common.FlagMap[peekt.raw.cgoType], peekt.raw.nano)
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
func CmdStack(inputPath string, outPath string, full bool) error {
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

	st := newStacker(inFile, outFile, full)
	st.analyze()

	inFile.Close()
	outFile.Close()
	return nil
}

