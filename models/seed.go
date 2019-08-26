package models

import "fmt"

//SeedDB seeds database with initial DEMO data, don't leave default admin password for production
func SeedDB() {
	//DEFAULT USERS ========================================
	admin := User{}
	if db.Where("role = ?", ADMIN).First(&admin); admin.ID == 0 {
		//no admins found, create a default one
		admin.FirstName = "Admin"
		admin.MiddleName = "A."
		admin.LastName = "Admin"
		admin.Role = ADMIN
		admin.Password = "admin"
		admin.Email = "admin@admin.org"
		if err := db.Create(&admin).Error; err != nil {
			panic(err)
		}
	}
	manager := User{}
	if db.Where("role = ?", MANAGER).First(&manager); manager.ID == 0 {
		//no managers found, create a default one
		manager.FirstName = "Manager"
		manager.MiddleName = "M."
		manager.LastName = "Manager"
		manager.Role = MANAGER
		manager.Password = "manager"
		manager.Email = "manager@manager.org"
		if err := db.Create(&manager).Error; err != nil {
			panic(err)
		}
	}
	customer := User{}
	if db.Where("role = ?", CUSTOMER).First(&customer); customer.ID == 0 {
		//no customer found, create a default one
		customer.FirstName = "Customer"
		customer.MiddleName = "C."
		customer.LastName = "Customer"
		customer.Role = CUSTOMER
		customer.Password = "customer"
		customer.Email = "customer@customer.org"
		if err := db.Create(&customer).Error; err != nil {
			panic(err)
		}
	}

	//DEFAULT SETTINGS ========================================
	if len(GetSetting("tel_1")) == 0 {
		tel1 := Setting{Code: "tel_1", Title: "Contact phone #1", Value: "+1 300 33-33-33"}
		tel2 := Setting{Code: "tel_2", Title: "Contact phone #2", Value: "+1 300 43-33-33"}
		tel3 := Setting{Code: "tel_3", Title: "Contact phone #3", Value: "+1 300 53-33-33"}
		if err := db.Create(&tel1).Error; err != nil {
			panic(err)
		}
		if err := db.Create(&tel2).Error; err != nil {
			panic(err)
		}
		if err := db.Create(&tel3).Error; err != nil {
			panic(err)
		}
	}
	if len(GetSetting("email")) == 0 {
		email := Setting{Code: "email", Title: "Contact email", Value: "email@example.org"}
		if err := db.Create(&email).Error; err != nil {
			panic(err)
		}
	}
	if len(GetSetting("order_email")) == 0 {
		email := Setting{Code: "order_email", Title: "Email for order notifications", Value: "order@example.org"}
		if err := db.Create(&email).Error; err != nil {
			panic(err)
		}
	}
	if len(GetSetting("title_suffix")) == 0 {
		suffix := Setting{Code: "title_suffix", Title: "Meta Title suffix", Value: "Awesome business skeletons"}
		if err := db.Create(&suffix).Error; err != nil {
			panic(err)
		}
	}

	//home page
	if len(GetSetting("home_id")) == 0 {
		//Description is page content ^_^ Published is false for home page, so it won't be shown as a separate page
		home := Page{
			Title: "Welcome To GinShop Skeleton Project",
			Description: `
			<p>Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.</p>
			<p>The standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested. Sections 1.10.32 and 1.10.33 from "de Finibus Bonorum et Malorum" by Cicero are also reproduced in their exact original form, accompanied by English versions from the 1914 translation by H. Rackham.</p>
			`,
			Published: false,
		}
		if err := db.Create(&home).Error; err != nil {
			panic(err)
		}

		homeID := Setting{Code: "home_id", Title: "Home Page ID", Value: fmt.Sprintf("%d", home.ID)}
		if err := db.Create(&homeID).Error; err != nil {
			panic(err)
		}
	}
	//about page
	if len(GetSetting("about_id")) == 0 {
		about := Page{Title: "About GinShop", Description: "GinShop is a skeleton project for basic e-shop web-sites. It demonstrates user authentication, 3 user roles: admin, manager & customer, product cart & checkout, email notifications, user dashboard for all roles, some ui animations via Animate On Scroll library, bootstrap 4 template, home page carousel and may be something more...", Published: true}
		if err := db.Create(&about).Error; err != nil {
			panic(err)
		}

		aboutID := Setting{Code: "about_id", Title: "About Page ID", Value: fmt.Sprintf("%d", about.ID)}
		if err := db.Create(&aboutID).Error; err != nil {
			panic(err)
		}
	}
	//DEFAULT CAROUSEL SLIDES ==============================
	var slides []Slide
	if db.Find(&slides); len(slides) == 0 {
		slide := Slide{Title: "Demo slide Title", NavigationURL: "/", FileURL: "/public/uploads/tables.jpg"}
		if err := db.Create(&slide).Error; err != nil {
			panic(err)
		}
		slide2 := Slide{Title: "Demo slide Title 2", NavigationURL: "/", FileURL: "/public/uploads/window_sills.jpg"}
		if err := db.Create(&slide2).Error; err != nil {
			panic(err)
		}
	}

	//DEFAULT CATEGORIES ====================================
	var categories []Category
	if db.Find(&categories); len(categories) == 0 {
		category := Category{Title: "Product Category 1", Description: "Category description goes here", Published: true}
		if err := db.Create(&category).Error; err != nil {
			panic(err)
		}
	}

	//DEFAULT PRODUCTS ====================================
	var products []Product
	if db.Find(&products); len(products) == 0 {
		category := Category{}
		db.First(&category)
		for i := 0; i < 20; i++ {
			product := Product{
				Title:       fmt.Sprintf("Demo product %d", i),
				CategoryID:  category.ID,
				Description: "Product description goes here",
				Published:   true,
				Recommended: true, //to show the product on the home page
				Images: []Image{
					Image{URL: "/public/uploads/apple.jpg"},
				},
			}
			if err := db.Create(&product).Error; err != nil {
				panic(err)
			}
		}
	}

	//DEFAULT MENUS ===========================================
	var menus []Menu
	if db.Find(&menus); len(menus) == 0 {
		main := Menu{Code: "main", Title: "Main menu"}
		if err := db.Create(&main).Error; err != nil {
			panic(err)
		}
		footer := Menu{Code: "footer", Title: "Footer menu"}
		if err := db.Create(&footer).Error; err != nil {
			panic(err)
		}
		//navbar menu items
		mi := MenuItem{MenuID: main.ID, Title: "Home", URL: "/", Ord: 1}
		if err := db.Create(&mi).Error; err != nil {
			panic(err)
		}
		about := Page{}
		db.First(&about, GetSetting("about_id"))
		mi = MenuItem{MenuID: main.ID, Title: "About", URL: about.URL(), Ord: 2}
		if err := db.Create(&mi).Error; err != nil {
			panic(err)
		}
		//footer menu items is shown in 4 columns, so it has by default 4 parent items
		for i := 0; i < 4; i++ {
			mi = MenuItem{
				MenuID: footer.ID,
				Title:  fmt.Sprintf("Lorem Ipsum %d", i+1),
				URL:    "#",
				Ord:    i + 1,
				Children: []MenuItem{
					MenuItem{MenuID: footer.ID, Title: fmt.Sprintf("Subitem %d", i+1), URL: "/", Ord: 1},
				},
			}
			if err := db.Create(&mi).Error; err != nil {
				panic(err)
			}
		}

	}
}
