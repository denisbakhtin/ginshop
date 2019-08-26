package models

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

//CUSTOMER is a customer role name
const CUSTOMER = "customer"

//MANAGER is a manager role name
const MANAGER = "manager"

//ADMIN is an admin role name
const ADMIN = "admin"

//ANONYMOUS is an anonymous role name
const ANONYMOUS = "anonymous"

//Login view model
type Login struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

//Register view model
type Register struct {
	FirstName       string `form:"first_name" binding:"required"`
	MiddleName      string `form:"middle_name" binding:"required"`
	LastName        string `form:"last_name" binding:"required"`
	Email           string `form:"email" binding:"required"`
	Password        string `form:"password" binding:"required"`
	PasswordConfirm string `form:"password_confirm" binding:"required"`
}

//Manage user view model
type Manage struct {
	FirstName       string `form:"first_name" binding:"required"`
	MiddleName      string `form:"middle_name" binding:"required"`
	LastName        string `form:"last_name" binding:"required"`
	Email           string `binding:"-"`
	Password        string `form:"password" binding:"required"`
	PasswordConfirm string `form:"password_confirm" binding:"required"`
}

//User type contains user info
type User struct {
	Model

	Email      string `form:"email" binding:"required"`
	FirstName  string `form:"first_name"`
	MiddleName string `form:"middle_name"`
	LastName   string `form:"last_name"`
	Password   string `form:"password" binding:"required"`
	Role       string `form:"role" binding:"required"`
}

//BeforeSave gorm hook
func (u *User) BeforeSave() (err error) {
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	u.Password = string(hash)
	return
}

//IsAdmin checks if user is admin
func (u *User) IsAdmin() bool {
	return u.Role == ADMIN
}

//IsCustomer checks if user is a customer
func (u *User) IsCustomer() bool {
	return u.Role == CUSTOMER
}

//IsManager checks if user is a manager
func (u *User) IsManager() bool {
	return u.Role == MANAGER
}
