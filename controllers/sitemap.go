package controllers

import (
	"fmt"
	"path"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/denisbakhtin/ginshop/config"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/denisbakhtin/sitemap"
)

//CreateXMLSitemap creates xml sitemap for search engines, and saves in public/sitemap folder
func CreateXMLSitemap() {
	logrus.Info("Starting XML sitemap generation")
	folder := path.Join(config.GetConfig().Public, "sitemap")
	domain := config.GetConfig().Domain
	now := time.Now()
	items := make([]sitemap.Item, 1)
	db := models.GetDB()

	//Home page
	items = append(items, sitemap.Item{
		Loc:        fmt.Sprintf("%s", domain),
		LastMod:    now,
		Changefreq: "daily",
		Priority:   1,
	})

	//Static pages
	var pages []models.Page
	db.Where("published = true").Find(&pages)
	for i := range pages {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s%s", domain, pages[i].URL()),
			LastMod:    pages[i].UpdatedAt,
			Changefreq: "monthly",
			Priority:   0.8,
		})
	}

	//Categories
	var categories []models.Category
	db.Where("published = true").Find(&categories)
	for i := range categories {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s%s", domain, categories[i].URL()),
			LastMod:    categories[i].UpdatedAt,
			Changefreq: "monthly",
			Priority:   0.8,
		})
	}

	//Static pages
	var products []models.Product
	db.Where("published = true").Find(&products)
	for i := range products {
		items = append(items, sitemap.Item{
			Loc:        fmt.Sprintf("%s%s", domain, products[i].URL()),
			LastMod:    products[i].UpdatedAt,
			Changefreq: "monthly",
			Priority:   0.8,
		})
	}

	if err := sitemap.SiteMap(path.Join(folder, "sitemap1.xml.gz"), items); err != nil {
		logrus.Error(err)
		return
	}
	if err := sitemap.SiteMapIndex(folder, "sitemap_index.xml", domain+"/public/sitemap/"); err != nil {
		logrus.Error(err)
		return
	}
	logrus.Info("XML sitemap has been generated in " + folder)
}
