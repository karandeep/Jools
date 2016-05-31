package main

import (
	"bytes"
	"log"
	"net/smtp"
)

func main() {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial("127.0.0.1:25")
	if err != nil {
		log.Fatal(err)
	}
	// Set the sender and recipient.
	c.Mail("sender@example.org")
	c.Rcpt("surajnalin@yahoo.com")
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Test email from Go!\n"
	msg := subject + mime + "<html><body><h1>Hello World!</h1></body></html>"
	buf := bytes.NewBufferString(msg)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
}
