package main

import "fmt"
import "log"

var version string = "undef"

func main() {
	// Start our runner.
	go runner()

	// Block until we're told to quit.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	log.Printf("Caught signal '%d', exiting.", s)
}
