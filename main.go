package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	flags "github.com/jessevdk/go-flags"
)

var version string = "undef"

type Options struct {
	Debug bool `env:"DEBUG"    long:"debug"    description:"enable debug"`
}

func main() {
	var opts Options

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.Debug {
		log.SetLevel(log.DebugLevel)
	}
}
