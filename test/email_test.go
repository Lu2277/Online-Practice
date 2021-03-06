package test

import (
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"testing"
)

func TestEmail(t *testing.T) {
	// 简单设置 log 参数
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	em := email.NewEmail()
	// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
	em.From = "UUU <2591134973@qq.com>"
	// 设置 receiver 接收方 的邮箱  此处也可以填写自己的邮箱， 就是自己发邮件给自己
	em.To = []string{"2591134973@qq.com"}
	// 设置主题
	em.Subject = "使用Golang发送邮件"
	// 简单设置文件发送的内容，暂时设置成纯文本
	em.Text = []byte("hello world， 咱们用 golang 发个邮件吧！！")
	//设置服务器相关的配置
	err := em.Send("smtp.qq.com:25", smtp.PlainAuth("", "2591134973@qq.com", "uhpitgbjrkirdjed", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("send successfully ... ")
}
