module github.com/mahmoudhossam/mail/cmd

go 1.16

require (
	github.com/emersion/go-imap v1.1.0
	github.com/mahmoudhossam/mail v0.0.0-20210725171241-5e8644a5afbb
	github.com/urfave/cli v1.22.5
)

replace github.com/mahmoudhossam/mail => ../mail
