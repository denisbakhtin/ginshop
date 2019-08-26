package models

//MenuItem type contains menu item info
type MenuItem struct {
	Model

	Title    string     `form:"title" binding:"required"`
	URL      string     `form:"url"`
	Ord      int        `form:"ord"`
	ParentID *uint64    `form:"parent_id"`
	Class    string     `form:"class"`
	MenuID   uint64     `form:"menu_id" binding:"required"`
	Children []MenuItem `gorm:"foreignkey:ParentID"`
}

//GetParent returns parent item
func (m *MenuItem) GetParent() MenuItem {
	parent := MenuItem{}
	if m.ParentID != nil {
		db.First(&parent, *m.ParentID)
	}
	return parent
}
