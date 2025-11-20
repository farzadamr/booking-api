package repository

import (
	"context"

	"github.com/farzadamr/booking-api/src/domain/model"
)

type BaseRepository[TEntity any] interface {
	Create(ctx context.Context, entity TEntity) (TEntity, error)
	Update(ctx context.Context, id int, entity map[string]interface{}) (TEntity, error)
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (TEntity, error)
	//TODO: implement GetByFilter() with pagination
}

type UserRepository interface {
	ExistsMobileNumber(ctx context.Context, mobileNumber string) (bool, error)
	FetchUserInfo(ctx context.Context, username string, password string) (model.User, error)
	GetDefaultRole(ctx context.Context) (roleId int, err error)
	CreateUser(ctx context.Context, u model.User) (model.User, error)
}
type RoleRepository interface {
	BaseRepository[model.Role]
}
