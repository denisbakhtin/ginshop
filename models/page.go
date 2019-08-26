package models

import (
	"fmt"
	"strings"
)

//Page type contains page info
type Page struct {
	Model

	Title           string `form:"title"`
	Description     string `form:"description"`
	Slug            string `form:"slug"`
	Published       bool   `form:"published"`
	MetaKeywords    string `form:"meta_keywords"`
	MetaDescription string `form:"meta_description"`
}

//URL returns article url
func (p *Page) URL() string {
	return fmt.Sprintf("/pages/%d-%s", p.ID, p.Slug)
}

//BeforeSave gorm hook
func (p *Page) BeforeSave() (err error) {
	if strings.TrimSpace(p.Slug) == "" {
		p.Slug = createSlug(p.Title)
	}
	return
}
