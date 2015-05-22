package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"os/user"
)

func main() {
	smtpAddr := "localhost:1025"

	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	username := "nobody"
	user, err := user.Current()
	if err == nil && user != nil && len(user.Username) > 0 {
		username = user.Username
	}

	fromAddr := username + "@" + host

	flag.StringVar(&smtpAddr, "smtp-addr", smtpAddr, "SMTP server address")
	flag.StringVar(&fromAddr, "from", fromAddr, "SMTP sender")

	flag.Parse()

	recip := flag.Args()
	if len(recip) == 0 {
		os.Exit(10)
	}

	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		os.Exit(11)
	}

	err = smtp.SendMail(smtpAddr, nil, fromAddr, recip, body)
	if err != nil {
		log.Fatal(err)
	}
}
