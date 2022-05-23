package test

import (
	"cloud-disk/core/define"
	"crypto/tls"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendMail(t *testing.T){
	e := email.NewEmail()
	e.From = "Get <1426887306@qq.com>"
	e.To = []string{"agarthasf@gmail.com"}

	e.Subject = "Cloud Disk Validation Code"
	e.HTML = []byte("<h1>Your Validation Code is 123456</h1>")
	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", "1426887306@qq.com",
		define.MailPassword, "smtp.qq.com"),
		&tls.Config{
			InsecureSkipVerify: true,
			ServerName: "smtp.qq.com",
		})
	if err != nil{
		t.Fatal(err)
	}
}