package db

import (
	"errors"

	"gorm.io/gorm"
)

// handleNotFound 处理查询结果，如果是 ErrRecordNotFound 则返回 nil, nil
func handleNotFound[T any](result *T, err error) (*T, error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

// FirstOrNil 查询第一条记录，如果不存在则返回 nil, nil（自动处理 ErrRecordNotFound）
// 使用示例：
//
//	user, err := FirstOrNil[model.User](db.WithContext(ctx).Where("id = ?", id))
func FirstOrNil[T any](tx *gorm.DB) (*T, error) {
	var dest T
	err := tx.First(&dest).Error
	return handleNotFound(&dest, err)
}

// TakeOrNil 查询任意一条记录，如果不存在则返回 nil, nil（自动处理 ErrRecordNotFound）
// 使用示例：
//
//	user, err := TakeOrNil[model.User](db.WithContext(ctx).Where("login_name = ?", loginName))
func TakeOrNil[T any](tx *gorm.DB) (*T, error) {
	var dest T
	err := tx.Take(&dest).Error
	return handleNotFound(&dest, err)
}

// LastOrNil 查询最后一条记录，如果不存在则返回 nil, nil（自动处理 ErrRecordNotFound）
// 使用示例：
//
//	user, err := LastOrNil[model.User](db.WithContext(ctx).Where("status = ?", status))
func LastOrNil[T any](tx *gorm.DB) (*T, error) {
	var dest T
	err := tx.Last(&dest).Error
	return handleNotFound(&dest, err)
}
