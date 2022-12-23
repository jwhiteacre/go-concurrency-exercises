//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
)

func handleSignal(p *MockProcess) {
	var called int64

	// Set up channel on which to send signal notifications.
	c := make(chan os.Signal, 1)

	// Notify on interrupt signals
	signal.Notify(c, os.Interrupt)

	for s := range c {
		fmt.Printf("\n%s\n", s)
		if called > 0 {
			os.Exit(1)
		}
		go p.Stop()
		called++
	}
}

func main() {

	// Create a process
	proc := MockProcess{}

	// set up signal handler
	go handleSignal(&proc)

	// Run the process (blocking)
	proc.Run()
}
