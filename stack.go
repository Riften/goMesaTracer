package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
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
	depth int
	duration int64
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

	writeBuf []*stackTrace
	bufSize int
}

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
	_, err = st.dest.Write([]byte(fmt.Sprintf("%d %s %d\n", tr.raw.count, FlagMap[tr.raw.cgoType], tr.duration)))
	if err != nil {
		fmt.Println("Error when write to stacker dest: ", err)
	}
}

func (st *stacker) analyze() {
	scanTrace(st.src, st.handleRawTrace)
}

func (st *stacker) Flush() {
	fmt.Println("Flush")
	var curTrace *stackTrace
	for ; st.bufSize>0; st.bufSize-=1 {
		curTrace = st.writeBuf[st.bufSize-1]
		fmt.Println("...", FlagMap[curTrace.raw.cgoType])
		st.writeTrace(curTrace)
	}
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
				st.writeBuf[st.bufSize] = peekt
				st.bufSize += 1

				if st.stack.Empty() {
					st.Flush()
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
	inFile, err := os.OpenFile(inputPath, os.O_RDONLY, os.ModePerm)
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

	outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Println("Error when open output stack file: ", err)
		return err
	}
	// do rw
	st := &stacker{
		dest:  outFile,
		src:   inFile,
		stack: lls.New(),
		writeBuf: make([]*stackTrace, 1000),
		bufSize: 0,
	}
	st.analyze()

	inFile.Close()
	outFile.Close()
	return nil
}

