package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/amartin3659/VacationHomeRental/internal/models"
)

func TestListenForMail(t *testing.T) {
  var processedData models.MailData
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
  testData := models.MailData{
    To: "test1@test.com",
    From: "test2@test.com",
    Subject: "",
    Content: "",
  }
  didRecieve := func(d models.MailData) {
    processedData = d 
  }

  listenForMail(didRecieve)
  mailChan <- testData

  time.Sleep(time.Millisecond * 10)

  if processedData != testData {
    t.Error("Expected data to match, but it did not")
  }
}

func TestSendMail(t *testing.T) {
	m := models.MailData{
    To:      "testemail@test.com",
		From:    "noreply@bungalow-bliss.com",
		Subject: "Receipt of a request for a reservation",
		Content: "",
	}
  var logBuf bytes.Buffer
  var errBuf bytes.Buffer
  app.ErrorLog.SetOutput(&errBuf) 
  app.InfoLog.SetOutput(&logBuf)
  defer func() {
    app.ErrorLog.SetOutput(os.Stdout)
    app.InfoLog.SetOutput(os.Stdout)
  }()
  sendMSG(m)
  logOutput := logBuf.String()
  errOutput := errBuf.String()
  fmt.Println("Log Output:", logOutput)
  fmt.Println("Error Output:", errOutput)

  if !strings.Contains(logOutput, "email sent out!") {
    t.Error("Error occured")
  }
}

func TestNoServerConnect(t *testing.T) {
	m := models.MailData{
    To:      "testemail@test.com",
		From:    "noreply@bungalow-bliss.com",
		Subject: "Receipt of a request for a reservation",
		Content: "",
	}
  var logBuf bytes.Buffer
  var errBuf bytes.Buffer
  app.ErrorLog.SetOutput(&errBuf) 
  app.InfoLog.SetOutput(&logBuf)
  defer func() {
    app.ErrorLog.SetOutput(os.Stdout)
    app.InfoLog.SetOutput(os.Stdout)
  }()
  sendMSG(m, "not localhost")
  logOutput := logBuf.String()
  errOutput := errBuf.String()
  fmt.Println("Log Output:", logOutput)
  fmt.Println("Error Output:", errOutput)

  if !strings.Contains(errOutput, "Did not connect") {
    t.Error("Error occured")
  }
}

func TestReadTemplate(t *testing.T) {
	m := models.MailData{
    To:      "testemail@test.com",
		From:    "noreply@bungalow-bliss.com",
		Subject: "Receipt of a request for a reservation",
		Content: "",
    Template: "basic",
	}
  var logBuf bytes.Buffer
  var errBuf bytes.Buffer
  app.ErrorLog.SetOutput(&errBuf) 
  app.InfoLog.SetOutput(&logBuf)
  defer func() {
    app.ErrorLog.SetOutput(os.Stdout)
    app.InfoLog.SetOutput(os.Stdout)
  }()
  sendMSG(m)
  logOutput := logBuf.String()
  errOutput := errBuf.String()
  fmt.Println("Log Output:", logOutput)
  fmt.Println("Error Output:", errOutput)

  if strings.Contains(errOutput, "Error reading file") {
    t.Error("Error occured")
  }
}

func TestNoTemplate(t *testing.T) {
	m := models.MailData{
    To:      "testemail@test.com",
		From:    "noreply@bungalow-bliss.com",
		Subject: "Receipt of a request for a reservation",
		Content: "",
    Template: "no-template",
	}
  var logBuf bytes.Buffer
  var errBuf bytes.Buffer
  app.ErrorLog.SetOutput(&errBuf) 
  app.InfoLog.SetOutput(&logBuf)
  defer func() {
    app.ErrorLog.SetOutput(os.Stdout)
    app.InfoLog.SetOutput(os.Stdout)
  }()
  sendMSG(m)
  logOutput := logBuf.String()
  errOutput := errBuf.String()
  fmt.Println("Log Output:", logOutput)
  fmt.Println("Error Output:", errOutput)

  if !strings.Contains(errOutput, "Error reading file") {
    t.Error("Error occured")
  }
}

func TestFailedToSend(t *testing.T) {
	m := models.MailData{
    To:      "",
		From:    "noreply@bungalow-bliss.com",
		Subject: "Receipt of a request for a reservation",
		Content: "",
	}
  var logBuf bytes.Buffer
  var errBuf bytes.Buffer
  app.ErrorLog.SetOutput(&errBuf) 
  app.InfoLog.SetOutput(&logBuf)
  defer func() {
    app.ErrorLog.SetOutput(os.Stdout)
    app.InfoLog.SetOutput(os.Stdout)
  }()
  sendMSG(m)
  logOutput := logBuf.String()
  errOutput := errBuf.String()
  fmt.Println("Log Output:", logOutput)
  fmt.Println("Error Output:", errOutput)

  if !strings.Contains(errOutput, "Could not send email") {
    t.Error("Error occured")
  }
}
