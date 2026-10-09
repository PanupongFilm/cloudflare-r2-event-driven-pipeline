package r2

import (
	"context"
	"server/database"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllUploadLog(ctx context.Context, limit, offset int) ([]database.UploadLog, int64, error)
	CreateUploadLog(ctx context.Context, data CreateUploadLogRequest) (*database.UploadLog, error)
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


func (r *repository) CreateUploadLog(ctx context.Context, data CreateUploadLogRequest) (*database.UploadLog, error) {
	
	// Convert DTO to DB Model
	log := &database.UploadLog{
		UserID:     data.UserID,
		AccountID:  data.AccountID,
		Action:     database.ActionType(data.Action),
		BucketName: data.BucketName,
		ObjectKey:  data.ObjectKey,
		ObjectSize: data.ObjectSize,
		ETag:       data.ETag,
		EventTime:  data.EventTime,
		Status:     database.StatusPending, 
	}

	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return nil, err
	}

	return log, nil
}
