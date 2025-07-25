package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func main() {
	startTime := time.Now()
	fmt.Println("timer start:", startTime)

	applyCmd := exec.Command("kubectl", "apply", "-f", "../policy/kubearmor/micro.yaml")
	if err := applyCmd.Run(); err != nil {
		fmt.Println(err)
		return
	}

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

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Resource: /usr/bin/ls") {
			endTime := time.Now()
			fmt.Printf("policy enforcement time: %v\n", endTime.Sub(startTime))
			cmd.Process.Kill()
			break
		}
	}

	applyCmd = exec.Command("kubectl", "delete", "-f", "../policy/kubearmor/micro.yaml")
	if err := applyCmd.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
