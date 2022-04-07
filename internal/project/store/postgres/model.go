package postgres

type projectModel struct {
	ID string `gorm:"type:uuid, primaryKey; column:id"`

	Name     string `gorm:"column: project_name"`
	BasePath string `gorm:"index; column:base_path"`
	RealPath string `gorm:"real_path"`
	IsSPA    bool   `gorm:"column:is_spa"`
}
