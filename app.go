package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func main() {

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
	fmt.Println("Start cli-test")

	// config
	//executeDuration := 10 * time.Second // in second
	breakDuration := 10 * time.Second // in second
	thereshold := 70.0                // Percentage

	// Run command
	cmd := exec.Command(cmdString)
	// Set stdout and stderr pipe
	// var stdoutBuf, stderrBuf bytes.Buffer
	// stdoutIn, _ := cmd.StdoutPipe()
	// stderrIn, _ := cmd.StderrPipe()
	// var errStdout, errStderr error
	// stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	// stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	// Start command
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	// Watch CPU and set thereshold
	cpuPercentDanger := make(chan float64, 1)
	// This is put in a goroutine, as the cpu percentage watcher. When the watcher hits the thereshold, it will go to the case to trigger the kill
	// Otherwise it will keep watching until the process is done
	go func(thereshold float64) {
		var cpuPercentage float64
		defer close(cpuPercentDanger)
		for {
			cpuPercentageList, _ := cpu.Percent(3*time.Second, false)
			cpuPercentage = cpuPercentageList[0]
			if cpuPercentage >= thereshold {
				break
			}
		}
		cpuPercentDanger <- cpuPercentage
		return
	}(thereshold)

	// Print stdout and stderr
	// go func() {
	// 	_, errStdout = io.Copy(stdout, stdoutIn)
	// }()
	// go func() {
	// 	_, errStderr = io.Copy(stderr, stderrIn)
	// }()

	// Wait for process to be done - non blocking
	processDone := make(chan error, 1)
	// This is put as goroutine to non-block the wait, creating race condition between 2 channels that will trigger the stop of "runCommand" function
	go func() {
		processDone <- cmd.Wait()
	}()

	select {
	//case <-time.After(executeDuration):
	case lastCPUPercentage := <-cpuPercentDanger:
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("failed to kill: ", err)
		}
		fmt.Printf("CPU Usage: %f . Restarting in 10 seconds\n", lastCPUPercentage)
		time.Sleep(breakDuration)
		return false, nil
	case err := <-processDone:
		if err != nil {
			err = fmt.Errorf("\nprocess done with error = %v", err)
			return true, err
		}
		log.Print("\nprocess done gracefully without error")
		return true, nil
	}
}
