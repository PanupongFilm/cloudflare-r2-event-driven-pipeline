package webhook

import (
	"context"
	"server/database"
	"server/internal/r2"
)

type Service interface {
	AddUploadLog(ctx context.Context, data r2.CreateUploadLogRequest) (*database.UploadLog, error)
}

type service struct{
	r2Repository r2.Repository
}

func NewService(r2Repository r2.Repository) Service{
	return &service{
		r2Repository: r2Repository,
	}
}


func (s *service) AddUploadLog(ctx context.Context, data r2.CreateUploadLogRequest) (*database.UploadLog, error){

	addUploadlog, err := s.r2Repository.CreateUploadLog(ctx,data)
	if err != nil{
		return nil, err
	}

	return addUploadlog, nil
}