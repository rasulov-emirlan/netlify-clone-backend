package postgres

import (
	"context"
	"errors"

	"github.com/rasulov-emirlan/netlify-clone-backend/internal/project"
	"gorm.io/gorm"
)

type Repository struct {
	conn *gorm.DB
}

func NewRepo(conn *gorm.DB) (*Repository, error) {
	if conn == nil {
		return nil, errors.New("project: connection to database can't be nil")
	}
	if err := conn.AutoMigrate(&projectModel{}); err != nil {
		return nil, err
	}
	return &Repository{
		conn: conn,
	}, nil
}

func (r *Repository) Create(ctx context.Context, p project.Project) (project.Project, error) {
	pm := projectModel{
		ID:       string(p.ID),
		Name:     p.Name,
		BasePath: p.BasePath,
		RealPath: p.RealPath,
		IsSPA:    p.IsSPA,
	}
	res := r.conn.Create(pm)
	if res.Error != nil {
		return project.Project{}, res.Error
	}
	if res.RowsAffected == 0 {
		return project.Project{}, errors.New("postgres: could't  insert into database, don't know why")
	}
	return p, nil
}

func (r *Repository) Read(ctx context.Context, id project.ID) (project.Project, error) {
	p := &projectModel{}
	res := r.conn.First(p, "id = ?", string(id))
	if res.Error != nil {
		return project.Project{}, res.Error
	}
	return project.Project{
		ID:       id,
		Name:     p.Name,
		BasePath: p.BasePath,
		RealPath: p.RealPath,
		IsSPA:    p.IsSPA,
	}, nil
}

func (r *Repository) Update(ctx context.Context, id project.ID, p project.Project) (project.Project, error) {
	panic("not imeplemented")
}
func (r *Repository) Delete(ctx context.Context, id project.ID) error {
	panic("not implemented")
}
