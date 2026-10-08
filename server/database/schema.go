package database

import (
	"gorm.io/gorm"
	"time"
)

// -------------------- User -------------------- 

type User struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserName  string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"user_name"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex" json:"email,omitempty"`
	Password  string         `gorm:"type:varchar(500);not null" json:"-"` // ไม่ส่งใน JSON
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	UploadLogs []UploadLog `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"upload_logs,omitempty"`
}

func (User) TableName() string {
	return "users"
}


// -------------------- Event Log From R2 -------------------- 

// ActionType enum for R2 actions
type ActionType string

const (
	ActionPutObject               ActionType = "PutObject"
	ActionCompleteMultipartUpload ActionType = "CompleteMultipartUpload"
	ActionDeleteObject            ActionType = "DeleteObject"
)

// StatusType enum for event status
type StatusType string

const (
	StatusPending   StatusType = "pending"
	StatusProcessed StatusType = "processed"
	StatusFailed    StatusType = "failed"
)

type UploadLog struct {
	ID          string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      *string        `gorm:"type:uuid;index" json:"user_id,omitempty"` // Foreign Key (nullable)
	AccountID   string         `gorm:"type:varchar(255);not null;index" json:"account_id"`
	Action      ActionType     `gorm:"type:varchar(50);not null;index" json:"action"`
	BucketName  string         `gorm:"type:varchar(255);not null;index" json:"bucket_name"`
	ObjectKey   string         `gorm:"type:varchar(1000);not null;index" json:"object_key"`
	ObjectSize  int64          `gorm:"type:bigint;not null" json:"object_size"`
	ETag        string         `gorm:"type:varchar(255);not null" json:"etag"`
	EventTime   time.Time      `gorm:"not null;index" json:"event_time"`
	Status      StatusType     `gorm:"type:varchar(50);default:'pending';index" json:"status"`
	ProcessedAt *time.Time     `json:"processed_at,omitempty"`
	FailReason  string         `gorm:"type:text" json:"fail_reason,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relation: UploadLog belongs to User
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
}

// TableName กำหนดชื่อ table
func (UploadLog) TableName() string {
	return "upload_logs"
}
