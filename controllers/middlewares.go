package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/denisbakhtin/ginshop/config"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//ContextData stores in gin context the common data, such as user info...
func ContextData() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if uID := session.Get(userIDKey); uID != nil {
			user := models.User{}
			models.GetDB().First(&user, uID)
			if user.ID != 0 {
				c.Set("User", &user)
			}
		}

		if config.GetConfig().SignupEnabled {
			c.Set("SignupEnabled", true)
		}
		c.Next()
	}
}

//AuthRequired grants access to authenticated users with a 'role' role
func AuthRequired(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		grant := false
		if user, _ := c.Get("User"); user != nil {
			u, ok := user.(*models.User)
			grant = ok && u.Role == role
		} else {
			grant = false
		}
		if grant {
			c.Next()
		} else {
			session := sessions.Default(c)
			session.AddFlash("Access forbidden, please sign in.")
			session.Save()
			c.Redirect(http.StatusFound, fmt.Sprintf("/signin?return=%s", url.QueryEscape(c.Request.RequestURI)))
			c.Abort()
		}
	}
}
