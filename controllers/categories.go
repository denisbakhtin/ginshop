package controllers

import (
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//CategoryGet handles GET /c/:slug route
func CategoryGet(c *gin.Context) {
	db := models.GetDB()
	category := models.Category{}

	idslug := c.Param("idslug")
	id := atouint(strings.Split(idslug, "-")[0])
	db.First(&category, id)
	db.Where("category_id in(?)", category.IDs()).Preload("Images").Find(&category.Products)
	if category.ID == 0 || !category.Published {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	//redirect to canonical url
	if c.Request.URL.Path != category.URL() {
		c.Redirect(303, category.URL())
		return
	}

	h := DefaultH(c)
	h["Title"] = category.Title
	h["Category"] = category
	h["Breadcrumbs"] = categoryBreadcrumbs(&category)
	h["MetaKeywords"] = category.MetaKeywords
	h["MetaDescription"] = category.MetaDescription
	c.HTML(http.StatusOK, "categories/show", h)
}

//CategoryIndex handles GET /admin/categories route
func CategoryIndex(c *gin.Context) {
	db := models.GetDB()
	var categories []models.Category
	db.Order("parent_id nulls first, ord asc").Find(&categories)
	h := DefaultH(c)
	h["Title"] = "Product categories"
	h["Categories"] = categories
	c.HTML(http.StatusOK, "categories/index", h)
}

//CategoryNew handles GET /admin/new_category route
func CategoryNew(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "New product category"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()

	c.HTML(http.StatusOK, "categories/form", h)
}

//CategoryCreate handles POST /admin/new_category route
func CategoryCreate(c *gin.Context) {
	db := models.GetDB()
	category := models.Category{}
	if err := c.ShouldBind(&category); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_category")
		return
	}
	if *category.ParentID == 0 {
		category.ParentID = nil
	}
	if err := db.Create(&category).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/categories")
}

//CategoryEdit handles GET /admin/categories/:id/edit route
func CategoryEdit(c *gin.Context) {
	db := models.GetDB()
	category := models.Category{}
	db.First(&category, c.Param("id"))
	if category.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = "Edit product category"
	h["Category"] = category
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "categories/form", h)
}

//CategoryUpdate handles POST /admin/categories/:id/edit route
func CategoryUpdate(c *gin.Context) {
	category := models.Category{}
	db := models.GetDB()
	if err := c.ShouldBind(&category); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/categories")
		return
	}
	if *category.ParentID == 0 {
		category.ParentID = nil
	}
	if err := db.Save(&category).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/categories")
}

//CategoryDelete handles POST /admin/categories/:id/delete route
func CategoryDelete(c *gin.Context) {
	category := models.Category{}
	db := models.GetDB()
	db.First(&category, c.Param("id"))
	if category.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	if err := db.Delete(&category).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/categories")
}
