package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginshop/config"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-gonic/gin"
	gomail "gopkg.in/gomail.v2"
)

//OrderConsultPost handles GET /orderconsult route
func OrderConsultPost(c *gin.Context) {
	orderconsult := models.OrderConsult{}
	c.ShouldBind(&orderconsult)

	//send email notification
	notifyAdminOfOrderConsult(c, &orderconsult)

	c.HTML(200, "orderconsult/confirm", gin.H{"OrderConsult": orderconsult})
}

func notifyAdminOfOrderConsult(c *gin.Context, orderconsult *models.OrderConsult) {
	//closure is needed here, as r may be released by the time func finishes
	go func() {
		var b bytes.Buffer

		domain := config.GetConfig().Domain
		tmpl := template.New("").Funcs(getFuncMap())
		workingdir, _ := os.Getwd()
		tmpl, _ = tmpl.ParseFiles(path.Join(workingdir, "views", "emails", "admin_orderconsult.gohtml"))
		if err := tmpl.Lookup("emails/admin_orderconsult").Execute(&b, gin.H{"OrderConsult": orderconsult}); err != nil {
			logrus.Error(err)
			return
		}

		smtp := config.GetConfig().SMTP
		msg := gomail.NewMessage()
		msg.SetHeader("From", smtp.From)
		msg.SetHeader("To", getSetting("order_email"))
		msg.SetHeader("Subject", fmt.Sprintf("Consultation order on %s", domain))
		msg.SetBody(
			"text/html",
			b.String(),
		)

		port, _ := strconv.Atoi(smtp.Port)
		dialer := gomail.NewPlainDialer(smtp.SMTP, port, smtp.User, smtp.Password)
		sender, err := dialer.Dial()
		if err != nil {
			logrus.Error(err)
			return
		}
		if err := gomail.Send(sender, msg); err != nil {
			logrus.Error(err)
			return
		}
	}()
}
