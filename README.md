mhsendmail
==========

A sendmail replacement which forwards mail to an SMTP server.

```bash
> go get github.com/mailhog/mhsendmail

> mhsendmail test@mailhog.local <<EOF
From: App <app@mailhog.local>
To: Test <test@mailhog.local>
Subject: Test message

Some content!
EOF
```

You can also set the from address (for SMTP `MAIL FROM`):

```bash
./mhsendmail --from="admin@mailhog.local" test@mailhog.local ...
```

Or pass in multiple recipients:

```bash
./mhsendmail --from="admin@mailhog.local" test@mailhog.local test2@mailhog.local ...
```

Or override the destination SMTP server:

```bash
./mhsendmail --smtp-addr="localhost:1026" test@mailhog.local ...
```

To use from php.ini

```
sendmail_path = /usr/local/bin/mhsendmail
```

### Licence

Copyright ©‎ 2015 - 2016, Ian Kent (http://iankent.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
