package controllers

import (
	"net/http"

	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-gonic/gin"
)

//HomeGet handles GET / route
func HomeGet(c *gin.Context) {
	db := models.GetDB()
	page := models.Page{}

	db.First(&page, getSetting("home_id"))
	h := DefaultH(c)
	h["Title"] = getSetting("title_suffix")
	h["Page"] = page

	c.HTML(http.StatusOK, "home/show", h)
}
