package controllers

import (
	"net/http"

	"github.com/denisbakhtin/ginshop/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

//CartType represents cart hash map, stored in session cookie
type CartType map[uint64]bool

//CartGet handles GET /cart route
func CartGet(c *gin.Context) {
	db := models.GetDB()
	var products []models.Product
	cart := getCart(c)
	db.Where("id in(?)", getCartProductIDs(cart)).Find(&products)
	h := DefaultH(c)
	h["Title"] = "Shopping cart"
	h["Products"] = products
	c.HTML(http.StatusOK, "cart/show", h)
}

//CartAdd handles POST /cart/add/:id route
func CartAdd(c *gin.Context) {
	session := sessions.Default(c)
	cart := getCart(c)
	id := atouint64(c.Param("id"))
	cart[id] = true
	session.Set("cart", cart)
	session.Save()
	c.JSON(200, len(cart))
}

//CartProcess handles POST /cart/process route
func CartProcess(c *gin.Context) {
	db := models.GetDB()
	cart := getCart(c)
	_ = cart
	_ = db
	//TODO: process order, send emails
	c.Redirect(http.StatusFound, "/")
}

//CartDelete handles POST /cart/delete/:id route
func CartDelete(c *gin.Context) {
	session := sessions.Default(c)
	cart := getCart(c)
	id := atouint64(c.Param("id"))
	delete(cart, id)
	session.Set("cart", cart)
	session.Save()
	c.Redirect(302, "/cart")
}

func getCart(c *gin.Context) CartType {
	session := sessions.Default(c)
	scart := session.Get("cart")
	cart := make(CartType)
	if scart != nil {
		cart = scart.(CartType)
	}
	return cart
}

func getCartProductIDs(cart CartType) []uint64 {
	i := 0
	keys := make([]uint64, len(cart))
	for k := range cart {
		if cart[k] {
			keys[i] = k
			i++
		}
	}
	return keys
}
