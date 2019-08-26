package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//SettingIndex handles GET /admin/settings route
func SettingIndex(c *gin.Context) {
	db := models.GetDB()
	var settings []models.Setting
	db.Order("id").Find(&settings)
	h := DefaultH(c)
	h["Title"] = "Settings"
	h["Settings"] = settings
	c.HTML(http.StatusOK, "settings/index", h)
}

//SettingNew handles GET /admin/new_setting route
func SettingNew(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "New Setting"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "settings/form", h)
}

//SettingCreate handles POST /admin/new_setting route
func SettingCreate(c *gin.Context) {
	db := models.GetDB()
	setting := models.Setting{}
	if err := c.ShouldBind(&setting); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_setting")
		return
	}

	if err := db.Create(&setting).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/settings")
}

//SettingEdit handles GET /admin/settings/:id/edit route
func SettingEdit(c *gin.Context) {
	db := models.GetDB()
	setting := models.Setting{}
	db.First(&setting, c.Param("id"))
	if setting.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = "Edit Setting"
	h["Setting"] = setting
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "settings/form", h)
}

//SettingUpdate handles POST /admin/settings/:id/edit route
func SettingUpdate(c *gin.Context) {
	setting := models.Setting{}
	db := models.GetDB()
	if err := c.ShouldBind(&setting); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/settings")
		return
	}
	if err := db.Save(&setting).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/settings")
}

//SettingDelete handles POST /admin/settings/:id/delete route
func SettingDelete(c *gin.Context) {
	setting := models.Setting{}
	db := models.GetDB()
	db.First(&setting, c.Param("id"))
	if setting.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	if err := db.Delete(&setting).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/settings")
}
