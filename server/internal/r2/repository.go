package r2

import (
	"context"
	"server/database"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllUploadLog(ctx context.Context, limit, offset int) ([]database.UploadLog, int64, error)
	
}


type repository struct {
	db *gorm.DB
}


func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}


func (r *repository) GetAllUploadLog(ctx context.Context, limit, offset int) ([]database.UploadLog, int64, error) {
	var logs []database.UploadLog
	var totalCount int64

	if err := r.db.WithContext(ctx).
		Model(&database.UploadLog{}).
		Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Preload("User").                             
		Order("created_at DESC").                     
		Limit(limit).                              
		Offset(offset).                             
		Find(&logs).Error                           

	if err != nil {
		return nil, 0, err
	}

	return logs, totalCount, nil
}
