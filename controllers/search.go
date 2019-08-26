package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-gonic/gin"
)

//SearchGet handles GET /search?search= route
func SearchGet(c *gin.Context) {
	db := models.GetDB()
	var products []models.Product

	db.Preload("Images").Where("lower(title) LIKE lower(?)", fmt.Sprintf("%%%s%%", strings.TrimSpace(c.Query("search")))).Find(&products)

	h := DefaultH(c)
	h["Title"] = fmt.Sprintf("Search results for: %s", c.Query("search"))
	h["Products"] = products
	c.HTML(http.StatusOK, "search/show", h)
}
