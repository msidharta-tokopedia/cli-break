package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {

	fmt.Println("Start cli-test")

	var finished bool
	for !finished {
		isFinished, err := runCommand("cli-test")
		if isFinished {
			finished = true
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}

func runCommand(cmdString string) (bool, error) {
	// config
	executeDuration := 10 * time.Second // in second
	breakDuration := 10 * time.Second   // in second

	// Run command
	cmd := exec.Command(cmdString)
	// Set stdout and stderr pipe
	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	// Start command
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	// Print stdout
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	// Set timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(executeDuration):
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill: ", err)
		}
		log.Println("\nprocess killed as timeout reached. Restarting in 10 seconds")
		time.Sleep(breakDuration)
		return false, nil
	case err := <-done:
		if err != nil {
			err = fmt.Errorf("\nprocess done with error = %v", err)
			return true, err
		}
		log.Print("\nprocess done gracefully without error")
		return true, nil
	}
}
