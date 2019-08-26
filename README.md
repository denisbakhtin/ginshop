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

# Installation

Clone repo or use `go get` to download the project, install `npm`, run `npm install`. For the first time run the executable with `seed=true` flag to seed the database with initial demo data. See `models/seed.go` for details. Build assets with `npm run build` or watch them with `npm run watch`.

# Make it your own

Let's say you want to create `Amazing Website`. Add a new `GitHub` repository `https://github.com/denisbakhtin/amazingshop` (replace that with your own).

Prepare `ginshop`: delete its `.git` directory.

Issue:

```
rm -rf src/github.com/denisbakhtin/ginshop/.git
```

Replace all references of `github.com/denisbakhtin/ginshop` with `github.com/denisbakhtin/amazingshop`:

```
grep -rl 'github.com/denisbakhtin/ginshop' ./ | xargs sed -i 's/github.com\/denisbakhtin\/ginshop/github.com\/denisbakhtin\/amazingshop/g'
```

Move all files to the new location:

```
mv src/github.com/denisbakhtin/ginshop/ src/github.com/denisbakhtin/amazingshop
```

And push it to the corresponding repo:

```
cd src/github.com/denisbakhtin/amazingshop
git init
git add --all .
git commit -m "Amazing Shop First Commit"
git remote add origin https://github.com/denisbakhtin/amazingshop.git
git push -u origin master
```

You can now go back to your `GOPATH` and check if everything is ok:

```
go install github.com/denisbakhtin/amazingshop

# Development
For watching go files and server restart you can use `https://github.com/pilu/fresh`. After installation run `fresh` in project directory.
```