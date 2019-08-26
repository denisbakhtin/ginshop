package models

import (
	"regexp"
	"strings"
	"time"

	"github.com/fiam/gounidecode/unidecode"
	"github.com/jinzhu/gorm"

	//postgres dialect, required by gorm
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Model is a tuned version of gorm.Model
type Model struct {
	ID        uint64     `form:"id" gorm:"primary_key"`
	CreatedAt time.Time  `binding:"-" form:"-"`
	UpdatedAt time.Time  `binding:"-" form:"-"`
	DeletedAt *time.Time `binding:"-" form:"-"`
}

var db *gorm.DB

//SetDB establishes connection to database and saves its handler into db *gorm.DB
func SetDB(connection string) {
	var err error
	db, err = gorm.Open("postgres", connection)
	if err != nil {
		panic(err)
	}
}

//GetDB returns database handler
func GetDB() *gorm.DB {
	return db
}

//AutoMigrate runs gorm auto migration
func AutoMigrate() {
	db.AutoMigrate(&User{}, &Page{}, &Menu{}, &MenuItem{}, &Category{}, &Product{}, &Image{}, &Setting{}, &Order{}, &Slide{})
}

//truncate truncates string to n runes
func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}

//createSlug makes url slug out of string
func createSlug(s string) string {
	s = strings.ToLower(unidecode.Unidecode(s))                     //transliterate if it is not in english
	s = regexp.MustCompile("[^a-z0-9\\s]+").ReplaceAllString(s, "") //spaces
	s = regexp.MustCompile("\\s+").ReplaceAllString(s, "-")         //spaces
	return s
}
