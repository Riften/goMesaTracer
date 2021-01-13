package stack

import (
	"fmt"
	"github.com/Riften/goMesaTracer/common"
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
	"os"
	"path"
	"strings"
	"time"
)

var hasDetail = false

func getCallTypeFromName(callName string) int {
	for cgoType, name := range common.FlagMap {
		if strings.TrimSuffix(name, "_BEGIN") == callName {
			return cgoType
		}
	}
	return -1
}

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
		if cgoType > common.Threshold {
			break
		}
		if c.count > 0 {
			st.writeLn(fmt.Sprintf("%s | %d | %s | %d ns | %f%%",
				strings.TrimSuffix(common.FlagMap[cgoType], "_BEGIN"),
				c.count,
				time.Duration(c.totalDuration).String(),
				c.totalDuration/int64(c.count),
				100*float64(c.totalDuration)/float64(st.stStatistic.endTime - st.stStatistic.startTime)))
		}
	}
	st.writeLn("")
}

func (st *stacker) comparison(){
	firstCallType := getCallTypeFromName(st.firstCallToCompare)
	secondCallType := getCallTypeFromName(st.secondCallToCompare)
	if firstCallType == -1 || secondCallType == -1 {
		fmt.Println("The calls to be compared not found!")
		return
	}
	st.writeLn(fmt.Sprintf("## Info for call comparison between %s and %s", st.firstCallToCompare, st.secondCallToCompare))
	st.writeLn(fmt.Sprintf("- Duration: %s / %s = %f",
		st.firstCallToCompare,
		st.secondCallToCompare,
		float64(st.stStatistic.calls[firstCallType].totalDuration)/float64(st.stStatistic.calls[secondCallType].totalDuration)))
	st.writeLn(fmt.Sprintf("- Count: %s / %s = %f",
		st.firstCallToCompare,
		st.secondCallToCompare,
		float64(st.stStatistic.calls[firstCallType].count)/float64(st.stStatistic.calls[secondCallType].count)))
	st.writeLn(fmt.Sprintf("- Duration: %s / Total = %f%%",
		st.firstCallToCompare,
		100*float64(st.stStatistic.calls[firstCallType].totalDuration)/float64(st.stStatistic.endTime - st.stStatistic.startTime)))
	st.writeLn("")

}

func (st *stacker) detail(){
	st.writeLn("## Call parameter detail infos")
	st.writeLn("Call Name | Count | Total Detail Number | Average Detail Number")
	st.writeLn("- | - | - | -")
	for cgoType, c :=range st.stStatistic.calls {
		if cgoType <= common.Threshold {
			continue
		}
		if c.count > 0 {
			st.writeLn(fmt.Sprintf("%s | %d | %d | %d",
				common.FlagMap[cgoType],
				c.count,
				c.totalDetail,
				c.totalDetail/int64(c.count)))
		}
	}
}

func (st *stacker) statisticRawTrace(r *rawTrace) {
	if r.cgoType > common.Threshold {
		hasDetail = true
		st.stStatistic.calls[r.cgoType].count += 1
		st.stStatistic.calls[r.cgoType].totalDetail += r.nano
		return
	}

	if r.nano > st.stStatistic.endTime {
		st.stStatistic.endTime = r.nano
	}

	if st.stStatistic.startTime == 0 {
		st.stStatistic.startTime = r.nano
	}

	if r.cgoType == 0 {
		// st.writeLn(FlagMap[0])
		// st.stStatistic.startTime = r.nano
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

func CmdStatistic(inputPath string, outPath string, callToCompare1 string, callToCompare2 string) error {
	inFile, err := os.OpenFile(inputPath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Error when open input trace file: ", err)
		return err
	}

	if outPath == "" {
		inputSuffix := path.Ext(inputPath)
		withoutSuffix := strings.TrimSuffix(inputPath, inputSuffix)
		outPath = withoutSuffix + "_statistic.md"
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
		stStatistic:         newStackStatistic(),
		firstCallToCompare:  callToCompare1,
		secondCallToCompare: callToCompare2,
	}

	st.statistic()

	// do comparison
	if callToCompare2 != "" {
		st.comparison()
	} else {
		fmt.Println("No second call to be compared. No comparison to do")
	}

	inFile.Close()
	outFile.Close()
	return nil
}