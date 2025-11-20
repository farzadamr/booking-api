package model

type User struct {
	BaseModel
	Name         string
	MobileNumber string
	Password     string
	Enabled      bool
	UserRoles    *[]UserRoles
}

type Role struct {
	BaseModel
	Name      string
	UserRoles *[]UserRoles
}

type UserRoles struct {
	BaseModel
	User   User
	Role   Role
	UserId int
	RoleId int
}
