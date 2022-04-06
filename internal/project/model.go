package project

type ID int64

type Project struct {
	ID   ID     `json:"id"`
	Name string `json:"projectName"`

	BasePath string `json:"basePath"`
	IsSPA    bool   `json:"isSPA"`

	ApplicationPort string `json:"applicationPort"`
}
