package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type filter struct {
	flags []bool
	dest io.Writer
	src io.Reader
}

func newFilterFromReader(src io.Reader) []bool {
	scanner := bufio.NewScanner(src)
	res := make([]bool, totalFlag)
	var nScan int
	var err error
	var cgoType int
	for scanner.Scan() {
		//tmpRawTrace := &rawTrace{}
		nScan, err = fmt.Sscanf(scanner.Text(), "%d", &cgoType)
		if err != nil {
			fmt.Println("Error when scan line: ", err)
		} else if nScan==0 {
			fmt.Println("Error: Sscanf parsed no param.")
			continue
		}
		res[cgoType] = true
		//handler(tmpRawTrace)
	}
	return res
}

func cmdFilter(inPath string, outPath string) error {
	inFile, err := os.OpenFile(inPath, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Error when open input trace file: ", err)
		return err
	}
	defer func() {
		err := inFile.Close()
		if err != nil {
			fmt.Println("Error when close input trace file: ", err)
		}
	}()


	return nil
}
