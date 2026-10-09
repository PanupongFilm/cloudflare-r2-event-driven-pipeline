package r2

import "time"

type CreateUploadLogRequest struct {
	UserID     *string `json:"user_id,omitempty"`
	AccountID  string  `json:"account_id" validate:"required"`
	Action     string  `json:"action" validate:"required,oneof=PutObject CompleteMultipartUpload DeleteObject"`
	BucketName string  `json:"bucket_name" validate:"required"`
	ObjectKey  string  `json:"object_key" validate:"required"`
	ObjectSize int64   `json:"object_size" validate:"required,min=0"`
	ETag       string  `json:"etag" validate:"required"`
	EventTime  time.Time `json:"event_time" validate:"required"`
}

