package main

import "C"
import (
	"fmt"
	"github.com/Riften/goMesaTracer/tracer"
)

// Initialize the routine used for record trace
func init() {
	fmt.Println("Initialize main routine")
	go func() {
		defer func() {
			fmt.Println("Routine End")
		}()
		fmt.Println("Routine start")
		var recved string
		var ok bool
		outCh := tracer.GlobalTracer.OutCh()
		for {
			recved, ok = <-outCh
			if !ok {
				fmt.Println("Tracer channel closed!")
				return
			} else {
				fmt.Println(recved)
			}
		}
	}()
}

func main() {

	return
}
