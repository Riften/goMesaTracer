package main

import "C"
import "fmt"

// Initialize the routine used for record trace
func init() {
	fmt.Println("Initialize main routine")
	go func() {
		defer func() {
			fmt.Println("Routine End")
		}()
		fmt.Println("Routine start")

	}()
}

func main() {

	return
}
