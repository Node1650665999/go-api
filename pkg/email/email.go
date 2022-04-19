package email

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"gin-api/pkg/config"
	"gin-api/pkg/time"
	"regexp"
)

//SendMail 发送邮件
func SendMail(subject, body string, to []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.GetString("mail.from"))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", formatBody(body))

	dialer := gomail.NewDialer(
		config.GetString("mail.host"),
		config.GetInt("mail.port"),
		config.GetString("mail.username"),
		config.GetString("mail.password"),
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: config.GetBool("mail.is_ssl")}
	return dialer.DialAndSend(m)
}

func formatBody(body string) string {
	re, _ := regexp.Compile(`\n`)
	body = fmt.Sprintf("发生时间 : %s <br> 错误信息 : %+v", time.CurrentDate(), body)
	return  re.ReplaceAllString(body, "<br>")
}
