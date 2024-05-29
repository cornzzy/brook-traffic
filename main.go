package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

type TrafficData struct {
	User  string `json:"user"`
	Bytes string `json:"bytes"`
}

func main() {
	file, err := os.Open("traffic.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	userBytes := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var data TrafficData
		line := scanner.Text()
		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			// if the line does not contain the expected keys, continue to the next line
			continue
		}
		// Convert the bytes string to an integer
		bytes, err := strconv.Atoi(data.Bytes)
		if err != nil {
			continue
		}
		userBytes[data.User] += bytes
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	// Directly marshal the userBytes map to the desired JSON format
	output, err := json.Marshal(userBytes)
	if err != nil {
		log.Fatalf("error marshalling result: %s", err)
	}

	pidFile, err := os.Open("pid")
	if err != nil {
		log.Fatalf("failed to open pid file: %s", err)
	}
	defer pidFile.Close()

	pidScanner := bufio.NewScanner(pidFile)
	if !pidScanner.Scan() {
		log.Fatal("pid file is empty")
	}
	pidStr := pidScanner.Text()
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Fatalf("failed to convert pid to integer: %s", err)
	}

	// Send SIGUSR1 signal to the process using shell command
	cmd := exec.Command("kill", "-SIGUSR1", strconv.Itoa(pid))
	err = cmd.Run()
	if err != nil {
		log.Fatalf("failed to send SIGUSR1 signal to process: %s", err)
	}

	fmt.Println(string(output))
}
