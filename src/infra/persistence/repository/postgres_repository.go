package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/farzadamr/booking-api/src/common"
	"github.com/farzadamr/booking-api/src/constant"
	"github.com/farzadamr/booking-api/src/pkg/service_errors"
	"gorm.io/gorm"
)

var SoftDeleteExp string = "id = ? and deleted_by is null"

type BaseRepository[TEntity any] struct {
	database *gorm.DB
}

func NewBaseRepository[TEntity any](db *gorm.DB) *BaseRepository[TEntity] {
	return &BaseRepository[TEntity]{database: db}
}

func (r BaseRepository[TEntity]) Create(ctx context.Context, entity TEntity) (TEntity, error) {
	tx := r.database.WithContext(ctx).Begin()
	err := tx.
		Create(&entity).
		Error
	if err != nil {
		tx.Rollback()
		return entity, err
	}
	tx.Commit()
	return entity, nil
}

func (r BaseRepository[TEntity]) Update(ctx context.Context, id int, entity map[string]interface{}) (TEntity, error) {
	snakeMap := map[string]interface{}{}
	for k, v := range entity {
		snakeMap[common.ToSnakeCase(k)] = v
	}
	snakeMap["modified_by"] = &sql.NullInt64{Int64: int64(ctx.Value(constant.UserIdKey).(float64)), Valid: true}
	snakeMap["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}

	model := new(TEntity)
	tx := r.database.WithContext(ctx).Begin()
	if err := tx.Model(model).
		Where(SoftDeleteExp, id).
		Updates(snakeMap).
		Error; err != nil {
		tx.Rollback()
		return *model, err
	}
	tx.Commit()
	return *model, nil
}

func (r BaseRepository[TEntity]) Delete(ctx context.Context, id int) error {
	tx := r.database.WithContext(ctx).Begin()

	model := new(TEntity)

	deleteMap := map[string]interface{}{
		"deleted_by": &sql.NullInt64{Int64: int64(ctx.Value(constant.UserIdKey).(float64)), Valid: true},
		"deleted_at": sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}

	if ctx.Value(constant.UserIdKey) == nil {
		return &service_errors.ServiceError{EndUserMessage: service_errors.PermissionDenied}
	}
	if cnt := tx.
		Model(model).
		Where(SoftDeleteExp, id).
		Updates(deleteMap).
		RowsAffected; cnt == 0 {
		tx.Rollback()
		return &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
	}
	tx.Commit()
	return nil
}

//TODO : implement pagination and GetById BaseRepository
