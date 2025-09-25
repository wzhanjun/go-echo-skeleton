package mail

import (
	_ "embed"
	"strconv"

	"github.com/wzhanjun/go-echo-skeleton/pkg/config"
	"gopkg.in/gomail.v2"
)

var (
	MailSubject = "Welcome"
)

func SendMail(mailTo []string, subject string, body string) error {
	mailConn := map[string]string{
		"user": config.Cfg.Email.User,
		"pass": config.Cfg.Email.Pass,
		"host": config.Cfg.Email.Host,
		"port": config.Cfg.Email.Port,
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(mailConn["user"], MailSubject))
	m.SetHeader("To", mailTo...)    //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}
