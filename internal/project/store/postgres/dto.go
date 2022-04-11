package postgres

import "github.com/rasulov-emirlan/netlify-clone-backend/internal/project"

func projectToService(p Project) (project.Project, error) {
	return project.Project{
		ID:             project.ID(p.ID),
		Name:           p.Name,
		CurrentVersion: p.CurrentVersion,
		BasePath:       p.BasePath,
		RealPath:       p.RealPath,
		AssetsRealPath: p.AssetsRealPath,
		IsSPA:          p.IsSPA,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}, nil
}

func projectsToService(p []Project) ([]project.Project, error) {
	res := make([]project.Project, len(p))
	for i, v := range p {
		res[i] = project.Project{
			ID:             project.ID(v.ID),
			Name:           v.Name,
			CurrentVersion: v.CurrentVersion,
			BasePath:       v.BasePath,
			RealPath:       v.RealPath,
			AssetsRealPath: v.AssetsRealPath,
			IsSPA:          v.IsSPA,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
		}
	}
	return res, nil
}
