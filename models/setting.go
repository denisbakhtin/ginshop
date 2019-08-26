package models

//Setting type contains settings info
type Setting struct {
	Model

	Code  string `binding:"required" form:"code"`
	Title string `form:"title"`
	Value string `form:"value"`
}

//GetSetting returns a site setting by its code, empty string otherwise
func GetSetting(code string) string {
	setting := Setting{}
	db.Where("code = ?", code).First(&setting)
	return setting.Value
}
