package config

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
)

type Config struct {
	PrintBody bool
	URL       *url.URL
	Tool      string
	EmailTo   string
	flags     *flag.FlagSet
	insecure  bool
	host      string
}

// New returns a new configuration object based on the passed in args
func New(args []string) *Config {
	var c = new(Config)

	// We create a new command-line flag parser instead of using the default.
	// Default package-level vars may or may not be set up how we want, but they
	// often present a small security risk.
	c.flags = flag.NewFlagSet(args[0], flag.ContinueOnError)

	c.flags.BoolVar(&c.insecure, "insecure", false, "set to true to use http instead of https")
	c.flags.StringVar(&c.host, "host", "localhost", "hostname for the search")
	c.flags.BoolVar(&c.PrintBody, "print", false, "print the web page's body to STDOUT")
	c.flags.StringVar(&c.Tool, "tool", "", "tool to run (oregonnews or libweb)")
	c.flags.StringVar(&c.EmailTo, "email-to", "", "who gets emails on outages")

	// flag.FlagSet force-writes to output, so to keep output sane and easy to
	// parse, we have to set up a fake
	c.flags.SetOutput(io.Discard)

	var err = c.flags.Parse(os.Args[1:])
	if err != nil {
		c.Usage(err)
	}
	if c.EmailTo == "" {
		c.Usage(fmt.Errorf("-email-to must be set"))
	}

	c.URL = new(url.URL)
	c.URL.Scheme = "https"
	c.URL.Host = c.host
	if c.insecure {
		c.URL.Scheme = "http"
	}

	log.Printf("INFO - using %q for search test", c.URL.String())

	return c
}

// Usage prints out the usage information and exits
func (c *Config) Usage(err error) {
	var code int

	if err != flag.ErrHelp {
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err)
		code = 2
	}

	fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]\n\n", c.flags.Name())
	// To make use of the built-in function, we must now set output again
	c.flags.SetOutput(os.Stderr)
	c.flags.PrintDefaults()

	os.Exit(code)
}
