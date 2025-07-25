package main

import (
	"fmt"
	"os/exec"
	"time"
)

func main() {
	for {
		cmd := exec.Command("/usr/bin/ls", "-l")
		_, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(1 * time.Millisecond) // 1ms wait
	}
}
