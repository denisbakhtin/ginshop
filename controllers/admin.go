package controllers

import (
	"github.com/gin-gonic/gin"
)

//AdminGet handles GET /admin route
func AdminGet(c *gin.Context) {
	c.Redirect(302, "/admin/products")
}
