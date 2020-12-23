package main

import (
	"fmt"
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
	"os"
	"path"
	"strings"
	"time"
)

func (st *stacker) statistic(){
	scanTrace(st.src, st.statisticRawTrace)
	st.writeLn("## Basic infos")
	st.writeLn(fmt.Sprintf("- Start: %s", time.Unix(0, st.stStatistic.startTime).Format("2006-01-02 15:04:05")))
	st.writeLn(fmt.Sprintf("- End: %s", time.Unix(0, st.stStatistic.endTime).Format("2006-01-02 15:04:05")))
	st.writeLn(fmt.Sprintf("- Total duration: %s", time.Duration(st.stStatistic.endTime - st.stStatistic.startTime).String()))
	st.writeLn(fmt.Sprintf("- Time extracted %s", time.Duration(st.stStatistic.busyTime).String()))
	st.writeLn("")
	st.writeLn("## Info for different call")
	st.writeLn("Flag Name | Count | Total Duration | Average Duration | Duration ratio")
	st.writeLn("- | - | - | - | -")
	for cgoType, c := range st.stStatistic.calls {
		if c.count > 0 {
			st.writeLn(fmt.Sprintf("%s | %d | %s | %d ns | %f%%",
				strings.TrimSuffix(FlagMap[cgoType], "_BEGIN"),
				c.count,
				time.Unix(0, c.totalDuration).String(),
				c.totalDuration/int64(c.count),
				100*float64(c.totalDuration)/float64(st.stStatistic.endTime - st.stStatistic.startTime)))
		}
	}
}

func (st *stacker) statisticRawTrace(r *rawTrace) {
	if r.nano > st.stStatistic.endTime {
		st.stStatistic.endTime = r.nano
	}

	if r.cgoType == 0 {
		//st.writeLn(FlagMap[0])
		st.stStatistic.startTime = r.nano
		return
	}
	if r.cgoType == 1 {
		//st.writeLn(FlagMap[1])
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

				st.stStatistic.calls[peekt.raw.cgoType].count += 1
				st.stStatistic.calls[peekt.raw.cgoType].totalDuration += peekt.duration
				if st.stack.Empty() {
					st.stStatistic.busyTime += peekt.duration
				}
				break
			}
		}
	}
}

func cmdStatistic(inputPath string, outPath string) error {
	inFile, err := os.OpenFile(inputPath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Error when open input trace file: ", err)
		return err
	}

	if outPath == "" {
		inputSuffix := path.Ext(inputPath)
		withoutSuffix := strings.TrimSuffix(inputPath, inputSuffix)
		outPath = withoutSuffix + "_statistic.txt"
		fmt.Println("No outPath specified. Out to default file ", outPath)
	}
	fmt.Println("Call statisthc will be writen to ", outPath)

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
		stStatistic: newStackStatistic(),
	}

	st.statistic()

	inFile.Close()
	outFile.Close()
	return nil
}