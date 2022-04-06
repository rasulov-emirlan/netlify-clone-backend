package project

type ID string

type Project struct {
	ID   ID     `json:"id"`
	Name string `json:"projectName"`

	// BasePath is the path from which we will
	// redirect our users to the RealPath
	BasePath string `json:"basePath"`

	// RealPath is the full url to the place
	// where all the files for this project will
	// be stored
	RealPath string `json:"realPath"`

	// If IsSPA is true than we will redirect all the
	// incomming requests for /BasePath/* to the index.html
	// at the root of your project
	IsSPA bool `json:"isSPA"`
}

func NewModel(name, basePath string, isSPA bool) (Project, error) {
	p := Project{
		Name:     name,
		BasePath: basePath,
		IsSPA:    isSPA,
	}
	return p, p.Validate()
}

func (p *Project) Validate() error {
	return nil
}
