package controllers

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/denisbakhtin/ginshop/config"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

const userIDKey = "UserID"

var tmpl *template.Template

//Option represents select option entry
type Option struct {
	Value string
	Text  string
}

//Breadcrumb represents a breadcrumb
type Breadcrumb struct {
	URL   string
	Title string
}

//DefaultH returns common to all pages template data
func DefaultH(c *gin.Context) gin.H {
	return gin.H{
		"Title":           "", //page title:w
		"Context":         c,
		"Csrf":            csrf.GetToken(c),
		"MetaKeywords":    "",
		"MetaDescription": "",
	}
}

//LoadTemplates loads templates from views directory to gin engine
func LoadTemplates(router *gin.Engine) {
	router.SetFuncMap(getFuncMap())
	router.LoadHTMLGlob("views/**/*")
}

func getFuncMap() template.FuncMap {
	return template.FuncMap{
		"isActiveLink":         isActiveLink,
		"stringInSlice":        stringInSlice,
		"formatDateTime":       formatDateTime,
		"now":                  now,
		"activeUserEmail":      activeUserEmail,
		"activeUserName":       activeUserName,
		"activeUserID":         activeUserID,
		"isUserAuthenticated":  isUserAuthenticated,
		"signUpEnabled":        SignUpEnabled,
		"noescape":             noescape,
		"footerMenuItems":      footerMenuItems,
		"topLevelMenuItems":    topLevelMenuItems,
		"refEqUint":            refEqUint,
		"selectableCategories": selectableCategories,
		"topLevelCategories":   topLevelCategories,
		"userRoles":            userRoles,
		"userRole":             userRole,
		"getSetting":           getSetting,
		"isNotBlank":           isNotBlank,
		"tel":                  tel,
		"productTitles":        productTitles,
		"cssVersion":           cssVersion,
		"jsVersion":            jsVersion,
		"domain":               domain,
		"isAdmin":              isAdmin,
		"isManager":            isManager,
		"isCustomer":           isCustomer,
		"slides":               homeSlides,
		"mainMenu":             mainMenu,
		"panelEntryPoint":      panelEntryPoint,
		"recommendedProducts":  recommendedProducts,
		"aboutURL":             aboutURL,
		"cartLen":              cartLength,
		"safeURL":              safeURL,
	}
}

//atouint converts string to uint, returns 0 if error
func atouint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 32)
	return uint(i)
}

//atouint64 converts string to uint64, returns 0 if error
func atouint64(s string) uint64 {
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}

//isActiveLink checks uri against currently active (uri, or nil) and returns "active" if they are equal
func isActiveLink(c *gin.Context, uri string) string {
	if c != nil && c.Request.RequestURI == uri {
		return "active"
	}
	return ""
}

//formatDateTime prints timestamp in human format
func formatDateTime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

//stringInSlice returns true if value is in list slice
func stringInSlice(value string, list []string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

//now returns current timestamp
func now() time.Time {
	return time.Now()
}

//activeUserEmail returns currently authenticated user email
func activeUserEmail(c *gin.Context) string {
	user := activeUser(c)
	if user != nil {
		return user.Email
	}
	return ""
}

//activeUserName returns currently authenticated user name
func activeUserName(c *gin.Context) string {
	user := activeUser(c)
	if user != nil {
		return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}
	return ""
}

//activeUserID returns currently authenticated user ID
func activeUserID(c *gin.Context) uint64 {
	user := activeUser(c)
	if user != nil {
		return user.ID
	}
	return 0
}

//activeUserRole returns currently authenticated user role
func activeUser(c *gin.Context) *models.User {
	if c != nil {
		u, _ := c.Get("User")
		if user, ok := u.(*models.User); ok {
			return user
		}
	}
	return nil
}

//isUserAuthenticated returns true is user is authenticated
func isUserAuthenticated(c *gin.Context) bool {
	user := activeUser(c)
	return user != nil
}

func isAdmin(c *gin.Context) bool {
	user := activeUser(c)
	return user != nil && user.IsAdmin()
}

func isManager(c *gin.Context) bool {
	user := activeUser(c)
	return user != nil && user.IsManager()
}

func isCustomer(c *gin.Context) bool {
	user := activeUser(c)
	return user != nil && user.IsCustomer()
}

//SignUpEnabled returns true if sign up is enabled by config
func SignUpEnabled() bool {
	return config.GetConfig().SignupEnabled
}

//noescape unescapes html content
func noescape(content string) template.HTML {
	return template.HTML(content)
}

//top level menu items
func topLevelMenuItems(menuID uint64) []models.MenuItem {
	db := models.GetDB()
	var menus []models.MenuItem
	db.Preload("Children").Where("parent_id is null and menu_id = ?", menuID).Order("ord asc").Find(&menus)
	return menus
}

//footer menu items
func footerMenuItems() []models.MenuItem {
	db := models.GetDB()
	menu := models.Menu{}
	db.Where("code=?", "footer").First(&menu)
	var menus []models.MenuItem
	db.Preload("Children").Where("menu_id = ?", menu.ID).Order("ord asc").Find(&menus, "parent_id is null")
	return menus
}

//refEqUint checks if *uint64 equals uint64
func refEqUint(ref *uint64, val uint64) bool {
	if ref == nil {
		return false
	}
	return *ref == val
}

//selectableCategories returns a slice of categories available for selection by products
func selectableCategories() []models.Category {
	db := models.GetDB()
	var categories []models.Category
	db.Where("NOT EXISTS(select 1 from categories as c where c.parent_id = categories.id)").Order("ord").Find(&categories)
	return categories
}

func topLevelCategories() []models.Category {
	db := models.GetDB()
	var categories []models.Category
	db.Preload("Children").Order("ord asc").Find(&categories, "parent_id is null")
	return categories
}

func mainMenu() []models.MenuItem {
	db := models.GetDB()
	menu := models.Menu{}
	db.Where("code=?", "main").First(&menu)
	var menus []models.MenuItem
	db.Preload("Children").Where("menu_id = ?", menu.ID).Order("ord asc").Find(&menus, "parent_id is null")
	return menus
}

func userRoles() []Option {
	return []Option{Option{Value: models.CUSTOMER, Text: "Customer"}, Option{Value: models.MANAGER, Text: "Manager"}, Option{Value: models.ADMIN, Text: "Administrator"}}
}

func userRole(role string) string {
	switch role {
	case models.CUSTOMER:
		return "Customer"
	case models.MANAGER:
		return "Manager"
	case models.ADMIN:
		return "Administrator"
	default:
		return "Unknown"
	}
}

//getSetting returns a site setting by its code or empty string
func getSetting(code string) string {
	return models.GetSetting(code)
}

func isNotBlank(content string) bool {
	return len(content) > 0 && content != "<p>&nbsp;</p>"
}

func tel(content string) string {
	reg, err := regexp.Compile("[^\\+0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(content, "")
	return processedString
}

func productTitles() []string {
	db := models.GetDB()
	var titles []string
	db.Model(&models.Product{}).Where("published = true").Order("title").Pluck("title", &titles)
	return titles
}

func cssVersion() string {
	return fileVersion(path.Join(config.GetConfig().Public, "assets", "main.css"))
}

func jsVersion() string {
	return fileVersion(path.Join(config.GetConfig().Public, "assets", "application.js"))
}

func fileVersion(path string) string {
	file, err := os.Stat(path)
	if err != nil {
		return timeToString(time.Now())
	}
	modified := file.ModTime()
	return timeToString(modified)
}

func timeToString(t time.Time) string {
	return fmt.Sprintf("%04d%02d%02d-%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

func domain() string {
	return config.GetConfig().Domain
}

//panelEntryPoint returns an entry point for authenticated users, same as panelEntryURL but with context
func panelEntryPoint(c *gin.Context) string {
	user := activeUser(c)
	url := "/"
	if user == nil {
		return url
	}
	switch user.Role {
	case models.ADMIN:
		url = "/admin/orders"
	case models.MANAGER:
		url = "/manager/orders"
	case models.CUSTOMER:
		url = "/customer/orders"
	default:
		url = "/"
	}
	return url
}

//panelEntryURL returns an entry point for authenticated users
func panelEntryURL(user models.User) string {
	url := "/"
	switch user.Role {
	case models.ADMIN:
		url = "/admin/orders"
	case models.MANAGER:
		url = "/manager/orders"
	case models.CUSTOMER:
		url = "/customer/orders"
	default:
		url = "/"
	}
	return url
}

func homeSlides() []models.Slide {
	var slides []models.Slide
	models.GetDB().Order("ord").Find(&slides)
	return slides
}

func recommendedProducts() []models.Product {
	var products []models.Product
	models.GetDB().
		Where("published = true and recommended = true").
		Order("id desc").Preload("Images").Limit(8).Find(&products)
	return products
}

func aboutURL() string {
	id := getSetting("about_id")
	page := models.Page{}
	models.GetDB().First(&page, id)
	return page.URL()
}

func cartLength(c *gin.Context) int {
	cart := getCart(c)
	return len(cart)
}

func categoryBreadcrumbs(c *models.Category) []Breadcrumb {
	return []Breadcrumb{
		Breadcrumb{Title: "Home", URL: "/"},
		Breadcrumb{Title: c.Title},
	}
}

func productBreadcrumbs(p *models.Product) []Breadcrumb {
	return []Breadcrumb{
		Breadcrumb{Title: "Home", URL: "/"},
		Breadcrumb{Title: p.Category.Title, URL: p.Category.URL()},
		Breadcrumb{Title: p.Title},
	}
}

func pageBreadcrumbs(p *models.Page) []Breadcrumb {
	return []Breadcrumb{
		Breadcrumb{Title: "Home", URL: "/"},
		Breadcrumb{Title: p.Title},
	}
}

func safeURL(url string) template.URL {
	return template.URL(url)
}
