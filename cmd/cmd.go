package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"os/user"
	"strings"

	"github.com/ogier/pflag"
)

// Go runs the MailHog sendmail replacement.
func Go() {
	smtpAddr := "localhost:1025"

	goflag := false
	for _, g := range os.Args[1:] {
		if strings.HasPrefix(g, "-") && !strings.HasPrefix(g, "--") {
			if strings.HasPrefix(g, "-from ") || strings.HasPrefix(g, "-from=") ||
				strings.HasPrefix(g, "-smtp-addr ") || strings.HasPrefix(g, "-smtp-addr=") {
				goflag = true
				break
			}
		}
	}

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
	var recip []string

	if goflag {
		flag.StringVar(&smtpAddr, "smtp-addr", smtpAddr, "SMTP server address")
		flag.StringVar(&fromAddr, "from", fromAddr, "SMTP sender")

		flag.Parse()
		recip = flag.Args()
	} else {
		pflag.StringVar(&smtpAddr, "smtp-addr", smtpAddr, "SMTP server address")
		pflag.StringVarP(&fromAddr, "from", "f", fromAddr, "SMTP sender")

		pflag.Parse()
		recip = pflag.Args()
	}

	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading stdin")
		os.Exit(11)
	}

	msg, err := mail.ReadMessage(bytes.NewReader(body))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing message body")
		os.Exit(11)
	}

	if len(recip) == 0 {
		// We only need to parse the message to get a recipient if none where
		// provided on the command line.
		recip = append(recip, msg.Header.Get("To"))
	}

	err = smtp.SendMail(smtpAddr, nil, fromAddr, recip, body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error sending mail")
		log.Fatal(err)
	}

}
