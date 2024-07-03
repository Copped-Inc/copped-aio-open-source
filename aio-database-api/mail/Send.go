package mail

import (
	"bytes"
	"fmt"
	"github.com/Copped-Inc/aio-types/console"
	"github.com/google/uuid"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"log"
)

func (m Mail) Send(to string) {

	t, _ := template.ParseFiles("mail/mail.html")
	var body bytes.Buffer

	err := t.Execute(&body, m)

	if err != nil {
		console.Log("Error", "Mail", "Send", err)
		return
	}

	server := mail.NewSMTPClient()
	server.Host = "mail.privateemail.com"
	server.Port = 587
	server.Username = "contact@copped-inc.com"
	server.Password = "TICRvplzwnUBLFuh"
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		console.Log("Error", "Mail", "Send", err)
		return
	}

	msgUUID, _ := uuid.NewRandom()
	msgID := fmt.Sprintf("<%s@mx.google.com>", msgUUID.String())

	email := mail.NewMSG()
	email.SetFrom("Copped AIO <contact@copped-inc.com>")
	email.AddTo(to)
	email.SetSubject(m.Title)
	email.AddHeader("Message-User", msgID)

	email.SetBody(mail.TextHTML, body.String())

	if email.Error != nil {
		log.Fatal(email.Error)
	}

	_ = email.Send(smtpClient)

}
