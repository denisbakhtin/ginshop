package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//MenuIndex handles GET /admin/menus route
func MenuIndex(c *gin.Context) {
	db := models.GetDB()

	var menus []models.Menu
	db.Order("id asc").Find(&menus)
	h := DefaultH(c)
	h["Title"] = "Menus"
	h["Menus"] = menus
	c.HTML(http.StatusOK, "menus/index", h)
}

//MenuNew handles GET /admin/new_menu route
func MenuNew(c *gin.Context) {
	h := DefaultH(c)

	h["Title"] = "New menu"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	h["Menu"] = models.Menu{}
	session.Save()

	c.HTML(http.StatusOK, "menus/form", h)
}

//MenuCreate handles POST /admin/new_menu route
func MenuCreate(c *gin.Context) {
	db := models.GetDB()
	menu := models.Menu{}
	if err := c.ShouldBind(&menu); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_menu")
		return
	}
	if err := db.Create(&menu).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/menus")
}

//MenuEdit handles GET /admin/menu/:id/edit route
func MenuEdit(c *gin.Context) {
	db := models.GetDB()
	menu := models.Menu{}
	menuID := c.Param("id")

	db.First(&menu, menuID)
	if menu.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = "Edit menu"
	h["Menu"] = menu
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "menus/form", h)
}

//MenuUpdate handles POST /admin/menu/:id/edit route
func MenuUpdate(c *gin.Context) {
	menu := models.Menu{}
	db := models.GetDB()

	if err := c.ShouldBind(&menu); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/menus")
		return
	}
	if err := db.Save(&menu).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/menus")
}

//MenuDelete handles POST /admin/menu/:id/delete route
func MenuDelete(c *gin.Context) {
	menu := models.Menu{}
	db := models.GetDB()
	menuID := c.Param("id")

	db.First(&menu, menuID)
	if menu.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	if err := db.Delete(&menu).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/menus")
}
