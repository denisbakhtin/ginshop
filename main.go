package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/claudiu/gocron"
	"github.com/denisbakhtin/ginshop/config"
	"github.com/denisbakhtin/ginshop/controllers"
	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func main() {
	gob.Register(controllers.CartType{})

	seed := flag.String("seed", "false", "Seed database: true, false")
	mode := flag.String("mode", "debug", "Application mode: debug, release, test")
	flag.Parse()

	gin.SetMode(*mode)

	initLogger()
	config.LoadConfig()
	models.SetDB(config.GetConnectionString())
	models.AutoMigrate()
	if *seed == "true" {
		models.SeedDB()
	}

	//Periodic tasks
	gocron.Every(1).Day().Do(controllers.CreateXMLSitemap)
	gocron.Start()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.StaticFS("/public", http.Dir(config.PublicPath())) //better use nginx to serve assets (Cache-Control, Etag, fast gzip, etc)
	controllers.LoadTemplates(router)                         //load views

	//setup sessions
	conf := config.GetConfig()
	store := cookie.NewStore([]byte(conf.SessionSecret))
	store.Options(sessions.Options{Path: "/", HttpOnly: true, MaxAge: 7 * 86400}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("ginshop-session", store))
	router.Use(controllers.ContextData())

	//setup csrf protection
	router.Use(csrf.Middleware(csrf.Options{
		Secret: conf.SessionSecret,
		ErrorFunc: func(c *gin.Context) {
			logrus.Error("CSRF token mismatch")
			controllers.ShowErrorPage(c, 400, fmt.Errorf("CSRF token mismatch"))
			c.Abort()
		},
	}))

	router.GET("/", controllers.HomeGet)
	router.NoRoute(controllers.NotFound)
	router.NoMethod(controllers.MethodNotAllowed)

	if config.GetConfig().SignupEnabled {
		router.GET("/signup", controllers.SignUpGet)
		router.POST("/signup", controllers.SignUpPost)
	}
	router.GET("/signin", controllers.SignInGet)
	router.POST("/signin", controllers.SignInPost)
	router.GET("/signout", controllers.SignoutGet)

	router.GET("/pages/:idslug", controllers.PageGet)
	router.GET("/c/:idslug", controllers.CategoryGet)
	router.GET("/p/:idslug", controllers.ProductGet)
	router.GET("/rss", controllers.RssGet)

	router.GET("/cart", controllers.CartGet)
	router.POST("/cart/add/:id", controllers.CartAdd)
	router.POST("/cart/delete/:id", controllers.CartDelete)

	router.GET("/new_order", controllers.OrderNew)
	router.POST("/new_order", controllers.OrderCreate)
	router.GET("/confirm_order/:id", controllers.OrderConfirm)

	router.GET("/search", controllers.SearchGet)

	router.POST("/orderconsult", controllers.OrderConsultPost)

	//admin area
	admin := router.Group("/admin")
	admin.Use(controllers.AuthRequired(models.ADMIN))
	{
		admin.POST("/upload", controllers.UploadPost) //image upload

		admin.GET("/users", controllers.UserIndex)
		admin.GET("/new_user", controllers.UserNew)
		admin.POST("/new_user", controllers.UserCreate)
		admin.GET("/users/:id/edit", controllers.UserEdit)
		admin.POST("/users/:id/edit", controllers.UserUpdate)
		admin.POST("/users/:id/delete", controllers.UserDelete)

		admin.GET("/pages", controllers.PageIndex)
		admin.GET("/new_page", controllers.PageNew)
		admin.POST("/new_page", controllers.PageCreate)
		admin.GET("/pages/:id/edit", controllers.PageEdit)
		admin.POST("/pages/:id/edit", controllers.PageUpdate)
		admin.POST("/pages/:id/delete", controllers.PageDelete)

		admin.GET("/menus", controllers.MenuIndex)
		admin.GET("/new_menu", controllers.MenuNew)
		admin.POST("/new_menu", controllers.MenuCreate)
		admin.GET("/menu/:id/edit", controllers.MenuEdit)
		admin.POST("/menu/:id/edit", controllers.MenuUpdate)
		admin.POST("/menu/:id/delete", controllers.MenuDelete)

		admin.GET("/menu/:id", controllers.MenuItemIndex)
		admin.GET("/menu/:id/new_item", controllers.MenuItemNew)
		admin.POST("/menu/:id/new_item", controllers.MenuItemCreate)
		admin.GET("/menu/:id/edit/:itemid", controllers.MenuItemEdit)
		admin.POST("/menu/:id/edit/:itemid", controllers.MenuItemUpdate)
		admin.POST("/menu/:id/delete/:itemid", controllers.MenuItemDelete)

		admin.GET("/categories", controllers.CategoryIndex)
		admin.GET("/new_category", controllers.CategoryNew)
		admin.POST("/new_category", controllers.CategoryCreate)
		admin.GET("/categories/:id/edit", controllers.CategoryEdit)
		admin.POST("/categories/:id/edit", controllers.CategoryUpdate)
		admin.POST("/categories/:id/delete", controllers.CategoryDelete)

		admin.GET("/products", controllers.ProductIndex)
		admin.GET("/new_product", controllers.ProductNew)
		admin.POST("/new_product", controllers.ProductCreate)
		admin.GET("/products/:id/edit", controllers.ProductEdit)
		admin.POST("/products/:id/edit", controllers.ProductUpdate)
		admin.POST("/products/:id/delete", controllers.ProductDelete)

		admin.POST("/new_image", controllers.ImageCreate)
		admin.POST("/images/:id/delete", controllers.ImageDelete)

		admin.GET("/settings", controllers.SettingIndex)
		admin.GET("/new_setting", controllers.SettingNew)
		admin.POST("/new_setting", controllers.SettingCreate)
		admin.GET("/settings/:id/edit", controllers.SettingEdit)
		admin.POST("/settings/:id/edit", controllers.SettingUpdate)
		admin.POST("/settings/:id/delete", controllers.SettingDelete)

		admin.GET("/orders", controllers.OrderIndex)
		admin.GET("/orders/:id", controllers.OrderGet)
		admin.POST("/orders/:id/delete", controllers.OrderDelete)

		admin.GET("/slides", controllers.SlideIndex)
		admin.GET("/new_slide", controllers.SlideNew)
		admin.POST("/new_slide", controllers.SlideCreate)
		admin.GET("/slides/:id/edit", controllers.SlideEdit)
		admin.POST("/slides/:id/edit", controllers.SlideUpdate)
		admin.POST("/slides/:id/delete", controllers.SlideDelete)
	}

	//manager area
	manager := router.Group("/manager")
	manager.Use(controllers.AuthRequired(models.MANAGER))
	{
		manager.GET("/orders", controllers.OrderIndex)
		manager.GET("/orders/:id", controllers.OrderGet)
		manager.GET("/manage", controllers.ManageGet)
		manager.POST("/manage", controllers.ManagePost)
	}

	//customer area
	customer := router.Group("/customer")
	customer.Use(controllers.AuthRequired(models.CUSTOMER))
	{
		customer.GET("/orders", controllers.OrderIndex)
		customer.GET("/orders/:id", controllers.OrderGet)
		customer.GET("/manage", controllers.ManageGet)
		customer.POST("/manage", controllers.ManagePost)
	}

	// Listen and server on 0.0.0.0:8081
	router.Run(":8085")
}

//initLogger initializes logrus logger with some defaults
func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
