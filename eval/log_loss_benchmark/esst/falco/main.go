package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getLogs() ([]string, error) {
	getPodCmd := exec.Command("kubectl", "get", "pods", "-n", "falco", "-l", "app.kubernetes.io/name=falco", "-o", "jsonpath={.items[0].metadata.name}")
	podNameBytes, err := getPodCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	podName := strings.TrimSpace(string(podNameBytes))
	fmt.Println("falco pod name:", podName)

	cmd := exec.Command("kubectl", "logs", "-n", "falco", podName)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, stderr.String())
	}

	logStr := strings.Split(out.String(), "\n")
	return logStr, nil
}

func main() {
	oldLogs, err := getLogs()
	if err != nil {
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
	newLogs, err := getLogs()
	if err != nil {
		fmt.Println(err)
		return
	}
	oldLogSet := make(map[string]bool)
	for _, log := range oldLogs {
		oldLogSet[log] = true
	}

	lsCount := 0
	for _, log := range newLogs {
		if !oldLogSet[log] && strings.Contains(log, "/usr/bin/ls") {
			lsCount++
		}
	}

	fmt.Printf("[%d]\n", lsCount)
}
