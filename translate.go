package main

import (
	"bufio"
	"fmt"
	"github.com/Riften/goMesaTracer/common"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

func scanAndTranslate(input io.Reader, output io.Writer) {
	scanner := bufio.NewScanner(input)
	var n_scan int
	var err error
	var count int
	var cgoType int
	var nano int64
	for scanner.Scan() {
		n_scan, err = fmt.Sscanf(scanner.Text(), "%d %d %d", &count, &cgoType, &nano)
		if err != nil {
			fmt.Println("Error when scan line: ", err)
		} else if n_scan==0 {
			fmt.Println("Error: Sscanf parsed no param.")
		} else {
			_, err = output.Write([]byte(fmt.Sprintf("%9d %s %s",
				count,
				time.Unix(nano/1e9, nano%1e9).Format("15:04:05.000000000"),
				common.FlagMap[cgoType])))
			if err != nil {
				fmt.Println("Error when write to translated file: ", err)
			}
		}
	}
}

func cmdTranslate(inputPath string, outPath string) error {
	inFile, err := os.OpenFile(inputPath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Error when open input trace file: ", err)
		return err
	}
	defer inFile.Close()

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
	defer outFile.Close()

	scanAndTranslate(inFile, outFile)
	return nil
}
