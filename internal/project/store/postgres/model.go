package postgres

type Project struct {
	ID string `gorm:"type:uuid;primary_key; column:id"`

	Name     string `gorm:"index:idx_member, unique; column: project_name"`
	BasePath string `gorm:"index:idx_member, unique; index; column:base_path"`
	RealPath string `gorm:"index:idx_member, unique; real_path"`
	IsSPA    bool   `gorm:"column:is_spa"`
}
