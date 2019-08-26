package controllers

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//SlideIndex handles GET /admin/slides route
func SlideIndex(c *gin.Context) {
	db := models.GetDB()
	var slides []models.Slide
	db.Order("ord").Find(&slides)
	h := DefaultH(c)
	h["Title"] = "List of slides"
	h["Slides"] = slides
	c.HTML(http.StatusOK, "slides/index", h)
}

//SlideNew handles GET /admin/new_slide route
func SlideNew(c *gin.Context) {
	h := DefaultH(c)
	h["Title"] = "New slide"
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	h["Slide"] = models.Slide{}
	session.Save()

	c.HTML(http.StatusOK, "slides/form", h)
}

//SlideCreate handles slide /admin/new_slide route
func SlideCreate(c *gin.Context) {
	db := models.GetDB()

	err := c.Request.ParseMultipartForm(32 << 20) // ~32MB
	if err != nil {
		logrus.Error(err)
		c.HTML(http.StatusBadRequest, "errors/400", nil)
		return
	}

	mpartFile, mpartHeader, err := c.Request.FormFile("upload")
	if err != nil {
		logrus.Error(err)
		c.HTML(http.StatusBadRequest, "errors/400", gin.H{"Error": err.Error()})
		return
	}
	defer mpartFile.Close()
	uri, err := saveFile(mpartHeader, mpartFile)
	if err != nil {
		logrus.Error(err)
		c.HTML(http.StatusInternalServerError, "errors/500", gin.H{"Error": err.Error()})
		return
	}

	slide := models.Slide{}
	if err := c.ShouldBind(&slide); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/new_slide")
		return
	}
	slide.FileURL = uri

	if err := db.Create(&slide).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", gin.H{"Error": err.Error()})
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/slides")
}

//SlideEdit handles GET /admin/slides/:id/edit route
func SlideEdit(c *gin.Context) {
	db := models.GetDB()
	slide := models.Slide{}
	db.First(&slide, c.Param("id"))
	if slide.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	h := DefaultH(c)
	h["Title"] = "Edit slide"
	h["Slide"] = slide
	session := sessions.Default(c)
	h["Flash"] = session.Flashes()
	session.Save()
	c.HTML(http.StatusOK, "slides/form", h)
}

//SlideUpdate handles slide /admin/slides/:id/edit route
func SlideUpdate(c *gin.Context) {
	slide := models.Slide{}
	db := models.GetDB()
	err := c.Request.ParseMultipartForm(32 << 20) // ~32MB
	if err != nil {
		logrus.Error(err)
		c.HTML(http.StatusBadRequest, "errors/400", nil)
		return
	}

	uri := ""
	mpartFile, mpartHeader, err := c.Request.FormFile("upload")
	if mpartFile != nil {
		defer mpartFile.Close()
		uri, _ = saveFile(mpartHeader, mpartFile)
	}

	if err := c.ShouldBind(&slide); err != nil {
		session := sessions.Default(c)
		session.AddFlash(err.Error())
		session.Save()
		c.Redirect(http.StatusSeeOther, "/admin/slides")
		return
	}
	if len(uri) > 0 {
		slide.FileURL = uri
	}

	if err := db.Save(&slide).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/slides")
}

//SlideDelete handles slide /admin/slides/:id/delete route
func SlideDelete(c *gin.Context) {
	slide := models.Slide{}
	db := models.GetDB()
	db.First(&slide, c.Param("id"))
	if slide.ID == 0 {
		c.HTML(http.StatusNotFound, "errors/404", nil)
		return
	}
	if err := db.Delete(&slide).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "errors/500", nil)
		logrus.Error(err)
		return
	}
	c.Redirect(http.StatusFound, "/admin/slides")
}
