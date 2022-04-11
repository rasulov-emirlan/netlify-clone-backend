package postgres

import (
	"time"
)

type Project struct {
	ID string `gorm:"type:uuid;primary_key; column:id"`

	Name           string `gorm:"index:idx_member, unique; column: project_name"`
	CurrentVersion uint16 `gorm:"index"`
	BasePath       string `gorm:"index:idx_member, unique; column:base_path"`
	RealPath       string `gorm:"index:idx_member, unique; real_path"`
	AssetsRealPath string `gorm:"index:idx_member, unique; assets_real_path"`
	IsSPA          bool   `gorm:"column:is_spa"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
