package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	cmd := exec.Command("karmor", "logs")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return
	}

	startTime := time.Now()
	applyLs := exec.Command("kubectl", "exec", "-it", "ubuntu", "--", "./for.sh", os.Args[1], os.Args[2])
	if err := applyLs.Run(); err != nil {
		fmt.Println(err)
		return
	}

	endTime := time.Now()

	fmt.Printf("Instruction execution time: %v\n", endTime.Sub(startTime))

	lsCount := 0
	logCh := make(chan string)
	timeout := time.After(20 * time.Second)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			logCh <- scanner.Text()
		}
		close(logCh)
	}()
	check := false
	for {
		select {
		case <-timeout:
			cmd.Process.Kill()
			check = true
			break
		case line, ok := <-logCh:
			if !ok {
				break
			}
			if strings.Contains(line, "Resource: /usr/bin/ls") {
				lsCount++
			}
		}
		if check {
			break
		}
	}

	fmt.Printf("[%d]\n", lsCount)
}
