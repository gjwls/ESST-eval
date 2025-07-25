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

	applyCmd := exec.Command("helm", "upgrade", "--namespace", "falco", "falco", "falcosecurity/falco", "--set", "tty=true", "-f", "../policy/falco/micro.yaml")
	if err := applyCmd.Run(); err != nil {
		fmt.Println(err)
		return
	}

	for {
		getPodCmd := exec.Command("kubectl", "get", "pods", "-n", "falco", "-l", "app.kubernetes.io/name=falco", "-o", "jsonpath={.items[0].metadata.name}")
		podNameBytes, err := getPodCmd.Output()
		if err != nil {
			fmt.Println(err)
			return
		}
		podName := strings.TrimSpace(string(podNameBytes))

		cmd := exec.Command("kubectl", "logs", "-n", "falco", podName)

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		}

		logStr := out.String()
		if strings.Contains(logStr, "/usr/bin/ls") {
			endTime := time.Now()
			fmt.Printf("policy enforcement time: %v\n", endTime.Sub(startTime))
			return
		}

		time.Sleep(1 * time.Millisecond)
	}
}
