package controllers

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//MenuItemIndex handles GET /admin/menu/:id route
func MenuItemIndex(c *gin.Context) {
	db := models.GetDB()
	menuID := c.Param("id")

	var items []models.MenuItem
	db.Where("menu_id = ?", menuID).Order("parent_id desc, ord asc").Find(&items)
	h := DefaultH(c)
	h["Title"] = "Menu Items"
	h["MenuID"] = menuID
	h["Items"] = items
	c.HTML(http.StatusOK, "menu_items/index", h)
}

//MenuItemNew handles GET /admin/menu/:id/new_item route
func MenuItemNew(c *gin.Context) {
	h := DefaultH(c)
	menuID := c.Param("id")

	h["Title"] = "New menu item"
	h["Item"] = models.MenuItem{MenuID: atouint64(menuID)}
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "menu_items/form", h)
}

//MenuItemCreate handles POST /admin/menu/:id/new_item route
func MenuItemCreate(c *gin.Context) {
	db := models.GetDB()
	menuID := c.Param("id")
	item := models.MenuItem{}
	if err := c.ShouldBind(&item); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/menu/%s/new_item", menuID))
		return
	}
	if *item.ParentID == 0 {
		item.ParentID = nil
	}
	if err := db.Create(&item).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/admin/menu/%s", menuID))
}

//MenuItemEdit handles GET /admin/menu/:id/edit/:itemid route
func MenuItemEdit(c *gin.Context) {
	db := models.GetDB()
	item := models.MenuItem{}
	itemID := c.Param("itemid")

	db.First(&item, itemID)
	if item.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = "Edit menu item"
	h["Item"] = item
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "menu_items/form", h)
}

//MenuItemUpdate handles POST /admin/menu/:id/edit/:itemid route
func MenuItemUpdate(c *gin.Context) {
	item := models.MenuItem{}
	db := models.GetDB()
	menuID := c.Param("id")
	itemID := c.Param("itemid")

	if err := c.ShouldBind(&item); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/menu/%s/edit/%s", menuID, itemID))
		return
	}
	if *item.ParentID == 0 {
		item.ParentID = nil
	}
	if err := db.Save(&item).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/admin/menu/%s", menuID))
}

//MenuItemDelete handles POST /admin/menu/:id/delete/:itemid route
func MenuItemDelete(c *gin.Context) {
	item := models.MenuItem{}
	db := models.GetDB()
	menuID := c.Param("id")
	itemID := c.Param("itemid")

	db.First(&item, itemID)
	if item.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	if err := db.Delete(&item).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/admin/menu/%s", menuID))
}
