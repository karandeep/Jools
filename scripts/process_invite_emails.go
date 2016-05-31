package main

import (
	"../consumer"
	"../model"
	"log"
	"strconv"
)

func main() {
	c, err := consumer.NewQueueConsumer("InviteEmails", strconv.Itoa(model.INVITE), "invite-emails-consumer")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Running forever")
	select {}

	if err := c.Shutdown(); err != nil {
		log.Fatalf("error during shutdown: %s", err)
	}
}
