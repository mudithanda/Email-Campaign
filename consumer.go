package main

import (
	"fmt"
	"log"
	"net/smtp"
	"sync"
	"time"
)

func emailWorker(id int, ch chan Recipient, wg *sync.WaitGroup) {
	defer wg.Done()

	smtpHost := "localhost"
	smtpPort := "1025"

	for recipient := range ch {

		// msgFormat := fmt.Sprintf("To: %s\r\nSubject: Test E-mail\r\n\r\n%s\r\n", recipient.Email, "Testing Email campaign...")
		// msg := []byte(msgFormat)

		msg, err := executeTemplate(recipient)
		if err != nil {
			fmt.Printf("Worker %d: Error parsing template for %s", id, recipient.Email)
			// add dead letter queue
			continue
		}

		fmt.Printf("Worker %d: Sending Email to%q\n", id, recipient.Email)

		// err := smtp.SendMail(smtpHost+":"+smtpPort, nil, "mhna@gmail.com", []string{recipient.Email}, msg)

		err = smtp.SendMail(smtpHost+":"+smtpPort, nil, "mhna@gmail.com", []string{recipient.Email}, []byte(msg))

		if err != nil {
			log.Printf("Worker %d: failed to send to %q: %v", id, recipient.Email, err)
			continue
		}

		time.Sleep(50 * time.Millisecond) // to avoid red limit (overwhelm server)

		fmt.Printf("Worker %d: Sent Email to %s \n", id, recipient.Email)
		//fmt.Println(id, recipient)
	}
}