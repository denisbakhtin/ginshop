package models

//OrderConsult type contains consult order info
type OrderConsult struct {
	Name  string `form:"order_name"`
	Phone string `form:"order_phone"`
}
