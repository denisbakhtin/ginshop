package models

//Slide type contains carousel slide info
type Slide struct {
	Model

	Title         string `form:"title" binding:"required"`
	Content       string `form:"content"`
	NavigationURL string `form:"navigation_url"`
	FileURL       string `form:"file_url"`
	Ord           int    `form:"ord"`
}
