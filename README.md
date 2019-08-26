# GinShop e-shop skeletop web-site

* Clean responsive design (screenshot)
* Product search with auto completion
* 3 hardcoded user groups, for simplicity: admins, managers, customers.
* Some nice css animations as well as animations on scroll via aos.js library.
* Admin dashboard with CKEditor 5
* Rss, xml sitemap
* Breadcrumbs, product image zoom & ajax add to cart request and more...

# Screenshots
## Home page
![](/public/images/carousel.jpg)
## Responsive
![](/public/images/responsive.jpg)
## Dashboard
![](/public/images/dashboard.jpg)
## Animations
![](/public/images/animations.jpg)
## Shopping cart
![](/public/images/shopping_cart.jpg)
## Custom 404, 405, 500 error pages
![](/public/images/error_page.jpg)

#Installation

Clone repo or use `go get` to download the project, install `npm`, run `npm install`. For the first time run the executable with `seed=true` flag to seed the database with initial demo data. See `models/seed.go` for details. Build assets with `npm run build` or watch them with `npm run watch`.