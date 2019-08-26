package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-gonic/gin"
)

//ImageCreate handles POST /admin/new_image route
func ImageCreate(c *gin.Context) {
	db := models.GetDB()

	err := c.Request.ParseMultipartForm(32 << 20) // ~32MB
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, "")
		return
	}

	mpartFile, mpartHeader, err := c.Request.FormFile("upload")
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}
	defer mpartFile.Close()
	uri, err := saveFile(mpartHeader, mpartFile)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		return
	}

	image := models.Image{URL: uri}

	if err := db.Create(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		logrus.Error(err)
		return
	}
	c.JSON(200, image)
}

//ImageDelete handles POST /admin/images/:id/delete route
func ImageDelete(c *gin.Context) {
	image := models.Image{}
	db := models.GetDB()
	db.First(&image, c.Param("id"))
	if image.ID == 0 {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	if err := db.Delete(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Error: err.Error()})
		logrus.Error(err)
		return
	}
	c.JSON(201, nil)
}
