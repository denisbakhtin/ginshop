package models

//Menu type contains menu info
type Menu struct {
	Model

	Code  string `form:"code" binding:"required"`
	Title string `form:"title" binding:"required"`
	Items []MenuItem
}
