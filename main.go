package main

import (
	"fmt"
	"log"
	"net/smtp"
	"om-gwtf/internal/config"
	"om-gwtf/internal/libweb"
	"om-gwtf/internal/oregonnews"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Get command-line flags
	var c = config.New(os.Args)

	// Pull env vars for sensitive data
	var host = os.Getenv("SMTP_HOST")
	var user = os.Getenv("SMTP_USERNAME")
	var pass = os.Getenv("SMTP_PASSWORD")
	var port, _ = strconv.Atoi(os.Getenv("SMTP_PORT"))

	var smtpErrs []string
	if host == "" {
		smtpErrs = append(smtpErrs, "SMTP_HOST cannot be blank")
	}
	if user == "" {
		smtpErrs = append(smtpErrs, "SMTP_USERNAME cannot be blank")
	}
	if pass == "" {
		smtpErrs = append(smtpErrs, "SMTP_PASSWORD cannot be blank")
	}
	if port < 1 {
		smtpErrs = append(smtpErrs, "SMTP_PORT must be a valid integer")
	}
	if len(smtpErrs) > 0 {
		c.Usage(fmt.Errorf("You must configure SMTP environment variables: %s", strings.Join(smtpErrs, ", ")))
	}

	var test func(*config.Config) bool
	switch c.Tool {
	case "oregonnews":
		test = oregonnews.Run
	case "libweb":
		test = libweb.Run
	default:
		c.Usage(fmt.Errorf(`-tool must be "oregonnews" or "libweb"`))
	}

	// TODO: here's where we need to add things like emailed alerts, customized
	// output options (say, in the email), maybe some extra-verbose logging
	// written somewhere on failures, etc.
	if test(c) {
		os.Exit(0)
	}

	// Go's SMTP package is a bit painful to use, but it's at least well-defined
	// and well-documented. And one of the key pieces: plain auth is *never* used
	// if a TLS connection isn't able to be made.
	var auth = smtp.PlainAuth("", user, pass, host)
	var lines = []string{
		"To: "+c.EmailTo,
		"From: "+user,
		"Subject: Site Outage Alert: "+c.Tool,
		"",
		"No text",
	}
	var to = strings.Split(c.EmailTo, ",")
	var msg = []byte(strings.Join(lines, "\r\n"))
	log.Printf("Sending message:%s", strings.Replace(string(msg), "\r\n", "\\r\\n", -1))
	var err = smtp.SendMail(fmt.Sprintf("%s:%d", host, port), auth, user, to, msg)
	if err != nil {
		log.Fatalf("Unable to send email: %s", err)
	}
	os.Exit(1)
}
