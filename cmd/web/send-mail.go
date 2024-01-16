package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/amartin3659/VacationHomeRental/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail(testFn func(models.MailData)) {
  go func() {
    for {
      msg := <-app.MailChan
      sendMSG(msg)
      testFn(msg)
    }
  }()
}

func sendMSG(m models.MailData, host ...string) {
  var hostString string
  if len(host) > 0 {
    hostString = host[0] 
  } else {
    hostString = "localhost"
  }
  server := mail.NewSMTPClient()
  server.Host = hostString
  server.Port = 1025
  server.KeepAlive = false
  server.ConnectTimeout = 10 * time.Second
  server.SendTimeout = 10 * time.Second

  client, err := server.Connect()
  if err != nil {
    errorLog.Println("Did not connect", err)
    return
  }

  email := mail.NewMSG()
  email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
  if m.Template == "" {
    email.SetBody(mail.TextHTML, m.Content)
  } else {
    data, err := os.ReadFile(fmt.Sprintf("./../../static/email/templates/%s.html", m.Template))
    if err != nil {
      errorLog.Println("Error reading file", err) 
    }
    mailTemplate := string(data)
    msgToSend := strings.Replace(mailTemplate, "[%E-MAIL-CONTENT%]", m.Content, 1)
    email.SetBody(mail.TextHTML, msgToSend)
  }

  err = email.Send(client)
  if err != nil {
    errorLog.Println("Could not send email", err)
  } else {
    infoLog.Println("email sent out!")
  }
}
