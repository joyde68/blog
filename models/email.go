package models

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

func SendEmail(co *Comment) {
	// 邮箱地址
	Smtp_Host := GetSetting("smtp-host")
	if Smtp_Host == "" {
		return
	}
	Smtp_Port := GetSetting("smtp-port")
	Smtp_User := GetSetting("smtp-user")
	Smtp_Password := GetSetting("smtp-password")
	nickname := GetSetting("smtp-nick-name")
	//to := []string{co.Email}
	to := co.Email
	auth := smtp.PlainAuth("", Smtp_User, Smtp_Password, Smtp_Host)

	subject := "标题 - Test Email"
	content_type := "Content-Type: text/html; charset=UTF-8"

	content := GetContentById(co.Cid)
	if content == nil {
		fmt.Println(co.Cid)
		return
	}
	data := map[string]interface{}{
		"link":      GetSetting("site_url"),
		"site":      GetSetting("site_title"),
		"author":    co.Author,
		"text":      template.HTML("回复" + co.Author + ":<br/>" + co.Content),
		"title":     content.Title,
		"permalink": GetSetting("site_url")+content.Slug,
	}

	t, err := template.New("email_template").Parse(GetSetting(`create-comment-template`))
	if co.Pid != 0 {
		pco := GetCommentById(co.Pid)
		data["author_p"] = pco.Author
		data["text_p"] = pco.Content
		t, err = template.New("email_template").Parse(GetSetting(`reply-comment-template`))
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	var bodyBuffer bytes.Buffer
	err = t.Execute(&bodyBuffer, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	body := bodyBuffer.String()

	//msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname + "<" + Smtp_User + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	msg := []byte("To: " + to + "\r\nFrom: " + nickname + "<" + Smtp_User + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err = smtp.SendMail(Smtp_Host+":"+Smtp_Port, auth, Smtp_User, []string{to}, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
