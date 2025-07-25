package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	startTime := time.Now()
	fmt.Println("timer start:", startTime)

	applyCmd := exec.Command("kubectl", "apply", "-f", "../policy/tracee/micro.yaml")
	if err := applyCmd.Run(); err != nil {
		fmt.Println(err)
		return
	}

	for {
		getPodCmd := exec.Command("kubectl", "get", "pods", "-n", "tracee", "-l", "app.kubernetes.io/name=tracee", "-o", "jsonpath={.items[0].metadata.name}")
		podNameBytes, err := getPodCmd.Output()
		if err != nil {
			fmt.Println(err)
			return
		}
		podName := strings.TrimSpace(string(podNameBytes))

		cmd := exec.Command("kubectl", "logs", "-n", "tracee", podName)

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		}

		logStr := out.String()
		if strings.Contains(logStr, "enforce-test-v5") {
			endTime := time.Now()
			fmt.Printf("%v\n", endTime.Sub(startTime))
			applyCmd = exec.Command("kubectl", "delete", "-f", "../policy/tracee/micro.yaml")
			if err := applyCmd.Run(); err != nil {
				fmt.Println(err)
				return
			}
			return
		}

		time.Sleep(1 * time.Millisecond)
	}
}
