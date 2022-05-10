package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rasulov-emirlan/netlify-clone-backend/internal/project"
	"gorm.io/gorm"
)

type repository struct {
	conn *gorm.DB
}

func NewRepo(conn *gorm.DB) (*repository, error) {
	if conn == nil {
		return nil, errors.New("project: connection to database can't be nil")
	}
	// TODO: code bellow causes panic investigate when deploying :)
	// if err := conn.Debug().AutoMigrate(Project{}); err != nil {
	// 	return nil, err
	// }
	return &repository{
		conn: conn,
	}, nil
}

func (r *repository) Create(ctx context.Context, p project.Project) (project.Project, error) {
	id := uuid.New().String()
	pm := Project{
		ID:       id,
		Name:     p.Name,
		BasePath: p.BasePath,
		RealPath: p.RealPath,
		IsSPA:    p.IsSPA,
	}
	res := r.conn.Create(&pm)
	if res.Error != nil {
		return project.Project{}, res.Error
	}
	if res.RowsAffected == 0 {
		return project.Project{}, errors.New("postgres: could't  insert into database, don't know why")
	}
	p.ID = project.ID(id)
	return p, nil
}

func (r *repository) Read(ctx context.Context, id project.ID) (project.Project, error) {
	p := Project{}
	res := r.conn.First(&p, "id = ?", string(id))
	if res.Error != nil {
		return project.Project{}, res.Error
	}
	return projectToService(p)
}

func (r *repository) List(ctx context.Context) ([]project.Project, error) {
	p := []Project{}
	res := r.conn.Find(&p)
	if res.Error != nil {
		return nil, res.Error
	}
	return projectsToService(p)
}

func (r *repository) Update(ctx context.Context, id project.ID, p project.Project) (project.Project, error) {
	pp := Project{
		ID:             string(id),
		Name:           p.Name,
		BasePath:       p.BasePath,
		RealPath:       p.RealPath,
		AssetsRealPath: p.AssetsRealPath,
		CurrentVersion: p.CurrentVersion,
		IsSPA:          p.IsSPA,
	}
	tx := r.conn.Save(&pp)
	if tx.Error != nil {
		return project.Project{}, tx.Error
	}
	return projectToService(pp)
}

func (r *repository) Delete(ctx context.Context, id project.ID) error {
	return r.conn.Delete(&Project{ID: string(id)}).Error
}
