package migration

import (
	"github.com/farzadamr/booking-api/src/constant"
	"github.com/farzadamr/booking-api/src/domain/model"
	"github.com/farzadamr/booking-api/src/infra/persistence/database"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

func Up1() {
	database := database.GetDB()

	createTables(database)
}

func createTables(database *gorm.DB) {
	tables := []interface{}{}

	tables = addNewTable(database, model.User{}, tables)
	tables = addNewTable(database, model.Role{}, tables)
	tables = addNewTable(database, model.UserRoles{}, tables)
	createDefaultUserInformation(database)
}

func addNewTable(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}

func createDefaultUserInformation(database *gorm.DB) {

	adminRole := model.Role{Name: constant.AdminRoleName}
	createRoleIfNotExists(database, &adminRole)

	defaultRole := model.Role{Name: constant.DefaultRoleName}
	createRoleIfNotExists(database, &defaultRole)

	u := model.User{Name: constant.DefaultUserName,
		MobileNumber: "09111112222"}
	pass := "12345678"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)

	createAdminUserIfNotExists(database, &u, adminRole.Id)

}

func createRoleIfNotExists(database *gorm.DB, r *model.Role) {
	exists := 0
	database.
		Model(&model.Role{}).
		Select("1").
		Where("name = ?", r.Name).
		First(&exists)
	if exists == 0 {
		database.Create(r)
	}
}

func createAdminUserIfNotExists(database *gorm.DB, u *model.User, roleId int) {
	exists := 0
	database.
		Model(&model.User{}).
		Select("1").
		Where("MobilePhone = ?", u.MobileNumber).
		First(&exists)
	if exists == 0 {
		database.Create(u)
		ur := model.UserRoles{UserId: u.Id, RoleId: roleId}
		database.Create(&ur)
	}
}
