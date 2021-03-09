package common

import (
	"bufio"
	"fmt"
	"os"
)

// FlagMap ...
var FlagMap []string // cgoFlag ==> Name of flag

func init() {
	fmt.Println("Initialize FlagMap from ", FlagListFile)
	FlagMap = make([]string, TotalFlag)
	input, err := os.OpenFile(FlagListFile, os.O_RDONLY, 0666)
	if err != nil {
		//panic(err)
		fmt.Println("Warning: Open flag list file failed: ", err)
		fmt.Println("Tracer may run without flag map")
		return
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
		} else if nScan < 2 {
			fmt.Println("Error: Sscanf parsed few param.")
		} else {
			FlagMap[cgoType] = cgoName
		}
	}
	//OutFlagMap()
}

// OutFlagMap ...
func OutFlagMap() {
	for cgoType, cgoName := range FlagMap {
		if cgoName != "" {
			fmt.Println(cgoType, "\t==>\t", cgoName)
		}
	}
}
