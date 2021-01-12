package common

import (
	"bufio"
	"fmt"
	"os"
)

var FlagMap []string // cgoFlag ==> Name of flag

func init() {
	fmt.Println("Initialize FlagMap from ", FlagListFile)
	input, err := os.OpenFile(FlagListFile, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer input.Close()

	/* Use scanner directly
	defineReg, err := regexp.Compile(defineFmt)
	if err != nil {
		panic(err)
	}
	 */

	var nScan int
	var cgoType int
	var cgoName string
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		// params := defineReg.FindStringSubmatch(scanner.Text())
		nScan, err = fmt.Sscanf(scanner.Text(), "#define %s %d", &cgoName, &cgoType)
		if err != nil {
			fmt.Println("Error when scan line: ", err)
		} else if nScan<2 {
			fmt.Println("Error: Sscanf parsed few param.")
		} else {
			FlagMap[cgoType] = cgoName
		}
	}
	outFlagMap()
}

func outFlagMap() {
	for cgoType, cgoName := range FlagMap {
		if cgoName != "" {
			fmt.Println(cgoType, "\t==>\t", cgoName)
		}
	}
}