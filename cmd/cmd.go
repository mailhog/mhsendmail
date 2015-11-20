package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"os/user"
	"regexp"
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

	if len(recip) == 0 {
		// We only need to parse the message to get a recipient if none where
		// provided on the command line.
//		re := regexp.MustCompile("(?im)^To: (.*)\r*\n$")
		re := regexp.MustCompile("(?im)^To: (.*)\r*$")
		n := bytes.IndexByte(body, 0)
		var bodyStr string;
		if n < 0  {
			bodyStr = string(body)
		} else {
			bodyStr = string(body[:n])
		}
		includedRecip := re.FindAllString(bodyStr, -1)
		if includedRecip == nil {
			fmt.Fprintln(os.Stderr, "missing recipient")
			os.Exit(10)
		}
		newRecip := make([]string, len(recip), len(recip)+len(includedRecip))
		copy(newRecip, recip)
		recip = append(newRecip, includedRecip...)
	}

	err = smtp.SendMail(smtpAddr, nil, fromAddr, recip, body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error sending mail")
		log.Fatal(err)
	}

}
