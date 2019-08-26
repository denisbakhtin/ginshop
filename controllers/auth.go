package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//SignInGet handles GET /signin route
func SignInGet(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "Sign in"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "auth/signin", h)
}

//SignInPost handles POST /signin route, authenticates user
func SignInPost(c *gin.Context) {
	session := sessions.Default(c)
	login := models.Login{}
	db := models.GetDB()
	returnURL := c.DefaultQuery("return", "/")
	if err := c.ShouldBind(&login); err != nil {
		session.AddFlash("Please provide correct data.")
		session.Save()
		c.Redirect(http.StatusFound, "/signin")
		return
	}

	user := models.User{}
	db.Where("email = lower(?)", login.Email).First(&user)

	if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)) != nil {
		logrus.Errorf("Login error, IP: %s, Email: %s", c.ClientIP(), login.Email)
		session.AddFlash("Email or password invalid")
		session.Save()
		c.Redirect(http.StatusFound, fmt.Sprintf("/signin?return=%s", url.QueryEscape(returnURL)))
		return
	}

	session.Set(userIDKey, user.ID)
	session.Save()
	c.Redirect(http.StatusFound, panelEntryURL(user))
}

//SignUpGet handles GET /signup route
func SignUpGet(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "Регистрация в системе"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "auth/signup", h)
}

//SignUpPost handles POST /signup route, creates new user
func SignUpPost(c *gin.Context) {
	session := sessions.Default(c)
	register := models.Register{}
	db := models.GetDB()
	if err := c.ShouldBind(&register); err != nil {
		session.AddFlash("Please provide correct data.")
		session.Save()
		c.Redirect(http.StatusFound, "/signup")
		return
	}
	if register.Password != register.PasswordConfirm {
		session.AddFlash("Password and password confirmation do not match.")
		session.Save()
		c.Redirect(http.StatusFound, "/signup")
		return
	}
	register.Email = strings.ToLower(register.Email)
	user := models.User{}
	db.Where("email = ?", register.Email).First(&user)
	if user.ID != 0 {
		session.AddFlash("User with this email already registered")
		session.Save()
		c.Redirect(http.StatusFound, "/signup")
		return
	}
	//create user
	user.Email = register.Email
	user.Password = register.Password
	user.FirstName = register.FirstName
	user.MiddleName = register.MiddleName
	user.LastName = register.LastName
	user.Role = models.CUSTOMER
	if err := db.Create(&user).Error; err != nil {
		session.AddFlash("Error whilst registering new user.")
		session.Save()
		logrus.Errorf("Error whilst registering user: %v", err)
		c.Redirect(http.StatusFound, "/signup")
		return
	}

	session.Set(userIDKey, user.ID)
	session.Save()
	c.Redirect(http.StatusFound, "/")
	return
}

//SignoutGet handles GET /signout route
func SignoutGet(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(userIDKey)
	session.Save()
	c.Redirect(http.StatusSeeOther, "/")
}

//ManageGet handles GET /:role/manage route
func ManageGet(c *gin.Context) {
	user := activeUser(c)
	h := DefaultH(c)
	h["Title"] = "Your account"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	manage := models.Manage{FirstName: user.FirstName, MiddleName: user.MiddleName, LastName: user.LastName, Email: user.Email}
	h["Manage"] = manage
	session.Save()

	tmpl := "auth/manage"
	switch user.Role {
	case models.MANAGER:
		tmpl = "auth/manager_manage"
	case models.CUSTOMER:
		tmpl = "auth/customer_manage"
	}
	c.HTML(http.StatusOK, tmpl, h)
}

//ManagePost handles POST /:role/manage route, updates user credentials
func ManagePost(c *gin.Context) {
	session := sessions.Default(c)
	db := models.GetDB()
	user := activeUser(c)
	manage := models.Manage{}
	url := c.Request.RequestURI

	if err := c.ShouldBind(&manage); err != nil {
		session.AddFlash("Please provide correct data.")
		session.Save()
		c.Redirect(http.StatusFound, url)
		return
	}
	if manage.Password != manage.PasswordConfirm {
		session.AddFlash("Password and password confirm do not match.")
		session.Save()
		c.Redirect(http.StatusFound, url)
		return
	}

	dbuser := models.User{}
	db.First(&dbuser, user.ID)

	if user.ID == 0 {
		logrus.Errorf("Account update error, IP: %s, Email: %s", c.ClientIP(), user.Email)
		session.AddFlash("Error while updating your account")
		session.Save()
		c.Redirect(http.StatusFound, url)
		return
	}

	dbuser.Password = manage.Password
	dbuser.FirstName = manage.FirstName
	dbuser.MiddleName = manage.MiddleName
	dbuser.LastName = manage.LastName

	if err := db.Save(&dbuser).Error; err != nil {
		logrus.Errorf("Account update error: %s", err.Error())
		session.AddFlash("Error while updating your account")
		session.Save()
		c.Redirect(http.StatusFound, url)
	}

	c.Redirect(http.StatusFound, panelEntryURL(*user))
}
