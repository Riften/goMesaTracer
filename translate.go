package main

import (
	"fmt"
	"github.com/Riften/goMesaTracer/stack"
	"os"
	"path"
	"strings"
)

func cmdTranslate(inputPath string, outPath string) error {
	inFile, err := os.OpenFile(inputPath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Error when open input trace file: ", err)
		return err
	}

	if outPath == "" {
		inputSuffix := path.Ext(inputPath)
		withoutSuffix := strings.TrimSuffix(inputPath, inputSuffix)
		outPath = withoutSuffix + "_trans.csv"
		fmt.Println("No outPath specified. Out to default file ", outPath)
	}
	fmt.Println("Translated trace will be writen to ", outPath)

	outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Error when open translated trace file: ", err)
		return err
	}

	stack.scanTrace(inFile, func(trace *stack.rawTrace) {
		_, err = outFile.WriteString(fmt.Sprintf("%d %s %d\n", trace.count, FlagMap[trace.cgoType], trace.nano))
		if err != nil {
			fmt.Println("Error when write to translated file: ", err)
		}
	})

	inFile.Close()
	outFile.Close()
	return nil
}
