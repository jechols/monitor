package main

import (
	"fmt"
	"om-gwtf/internal/config"
	"om-gwtf/internal/oregonnews"
	"os"
)

func main() {
	var test func(*config.Config) bool
	var c = config.New(os.Args)
	switch c.Tool {
	case "oregonnews":
		test = oregonnews.Run
	case "libweb":
	default:
		c.Usage(fmt.Errorf(`-tool must be "oregonnews" or "libweb"`))
	}

	// TODO: here's where we need to add things like emailed alerts, customized
	// output options (say, in the email), maybe some extra-verbose logging
	// written somewhere on failures, etc.
	if test(c) {
		os.Exit(0)
	}
	os.Exit(1)
}
