package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println("Debugger $>")
	firstArg := os.Args[1]
	cmd := exec.Command(firstArg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// use ptrace on the child process
	cmd.SysProcAttr = &unix.SysProcAttr{
		Ptrace: true,
	}

	// start new process which is traced so it stops before
	// executing first instruction and sends signal to the parent process
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// wait for the target process to return a signal
	// SIGTRAP, the breakpoint trap signal
	cmd.Wait()
	//process is restarted and parent waits for its termination
	pid := cmd.Process.Pid
	fmt.Println("Restarting the Process", pid)
	// target process to run to completion,
	if err := unix.PtraceCont(pid, 0); err != nil {
		log.Fatal(err)
	}
}
