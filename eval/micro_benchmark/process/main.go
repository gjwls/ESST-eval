package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	numIterations := 10000

	for j := 0; j < 10; j++ {
		start := time.Now()

		for i := 0; i < numIterations; i++ {
			cmd := exec.Command("/usr/bin/ls")
			if err := cmd.Run(); err != nil {
				fmt.Println(err)
				return
			}
		}

		duration := time.Since(start)
		fmt.Printf("%v\n", duration)

	}

}
