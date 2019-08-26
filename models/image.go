package models

//Image type contains image info
type Image struct {
	Model

	URL       string  `form:"url"`
	ProductID uint64  `form:"product_id"`
	Hash      string  `gorm:"-" form:"-"`
	Product   Product `gorm:"save_associations:false" binding:"-" form:"-"`
}
