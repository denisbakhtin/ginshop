package controllers

import (
	"fmt"
	"time"

	"github.com/denisbakhtin/ginshop/config"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"
)

//RssGet handles GET /rss route
func RssGet(c *gin.Context) {
	now := time.Now()
	domain := config.GetConfig().Domain
	db := models.GetDB()

	feed := &feeds.Feed{
		Title:       "GinShop skeleton",
		Link:        &feeds.Link{Href: domain},
		Description: "A gin-powered e-shop boilerplate",
		Author:      &feeds.Author{Name: "GinShop Inc.", Email: getSetting("contact_email")},
		Created:     now,
		Copyright:   "All rights reserved © GinShop",
	}

	feed.Items = make([]*feeds.Item, 0)
	var pages []models.Page
	db.Where("published = true").Find(&pages)

	for i := range pages {
		feed.Items = append(feed.Items, &feeds.Item{
			Id:      fmt.Sprintf("%s/pages/%d", domain, pages[i].ID),
			Title:   pages[i].Title,
			Link:    &feeds.Link{Href: fmt.Sprintf("%s/pages/%d", domain, pages[i].ID)},
			Author:  &feeds.Author{Name: "GinShop Inc."},
			Created: now,
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		c.HTML(500, "errors/500", nil)
		return
	}
	c.Data(200, "text/xml", []byte(rss))
}
