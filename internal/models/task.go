package models

import (
	"time"

	"gorm.io/gorm"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	ID          string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Priority    Priority       `gorm:"type:varchar(10);not null" json:"priority"`
	ReminderAt  time.Time      `json:"reminder_at"`
	Done        bool           `gorm:"default:false" json:"done"`
	ParentID    *string        `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	Parent      *Task          `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"parent,omitempty"`
	Children    []Task         `gorm:"foreignKey:ParentID;references:ID" json:"children,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
