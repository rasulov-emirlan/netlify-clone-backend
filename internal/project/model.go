package project

import "time"

type ID string

type Project struct {
	ID   ID     `json:"id"`
	Name string `json:"projectName"`
	// We assume that we cannot delete previous versions
	// so if we are in version n then we assume that we have n-1 versions
	// of backups
	CurrentVersion uint16 `json:"currVersion"`

	// BasePath is the path from which we will
	// redirect our users to the RealPath
	BasePath string `json:"basePath"`

	// RealPath is the full url to the place
	// where all the files for this project will
	// be stored
	RealPath string `json:"realPath"`

	// All files that do not have following extensions 'html', 'js', 'css',
	// will be considered to be assets. And they are not stored as backups
	// so they have a different realPath
	AssetsRealPath string `json:"assetsRealPath"`

	// If IsSPA is true than we will redirect all the
	// incomming requests for /BasePath/* to the index.html
	// at the root of your project
	IsSPA bool `json:"isSPA"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewModel(name, basePath string, isSPA bool) (Project, error) {
	p := Project{
		Name:           name,
		BasePath:       basePath,
		CurrentVersion: 1,
		IsSPA:          isSPA,
	}
	return p, p.Validate()
}

func (p *Project) Validate() error {
	return nil
}
