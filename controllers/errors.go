package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//NotFound handles gin NotFound error
func NotFound(c *gin.Context) {
	ShowErrorPage(c, http.StatusNotFound, nil)
}

//MethodNotAllowed handles gin MethodNotAllowed error
func MethodNotAllowed(c *gin.Context) {
	ShowErrorPage(c, http.StatusMethodNotAllowed, nil)
}

//AccessForbidden handles Access Forbidden http code
func AccessForbidden(c *gin.Context) {
	ShowErrorPage(c, http.StatusForbidden, nil)
}

//InternalError handles Internal Server Error http code
func InternalError(c *gin.Context) {
	ShowErrorPage(c, http.StatusInternalServerError, nil)
}

//ShowErrorPage executes error template given its code
func ShowErrorPage(c *gin.Context, code int, err error) {
	H := DefaultH(c)
	H["Error"] = err
	H["Title"] = fmt.Sprintf("%d Error", code)
	c.HTML(code, fmt.Sprintf("errors/%d", code), H)
}
