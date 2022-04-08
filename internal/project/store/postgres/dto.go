package postgres

import "github.com/rasulov-emirlan/netlify-clone-backend/internal/project"

func projectToService(p Project) (project.Project, error) {
	return project.Project{
		ID:       project.ID(p.ID),
		Name:     p.Name,
		BasePath: p.BasePath,
		RealPath: p.RealPath,
		IsSPA:    p.IsSPA,
	}, nil
}

func projectsToService(p []Project) ([]project.Project, error) {
	res := make([]project.Project, len(p))
	for i, v := range p {
		res[i] = project.Project{
			ID:       project.ID(v.ID),
			Name:     v.Name,
			BasePath: v.BasePath,
			RealPath: v.RealPath,
			IsSPA:    v.IsSPA,
		}
	}
	return res, nil
}
