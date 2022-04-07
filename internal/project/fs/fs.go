package fs

import (
	"context"
	"io"
	"mime/multipart"
	"os"
)

type fs struct {
	dirname string
}

func NewFileSystem(dirname string) (*fs, error) {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		if err = os.Mkdir(dirname, os.ModePerm); err != nil {
			return nil, err
		}
	}
	return &fs{}, nil
}

func (f *fs) Upload(ctx context.Context, files []*multipart.FileHeader, id string) (string, error) {
	if _, err := os.Stat(f.dirname + "/" + id); os.IsExist(err) {
		return "", err
	}
	// TODO: solve issue with permissions
	if err := os.Mkdir(f.dirname+"/"+id, 0777); err != nil {
		return "", err
	}

	for i := range files { // loop through the files one by one
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			return "", err
		}

		out, err := os.Create(f.dirname + "/" + files[i].Filename)

		defer out.Close()
		if err != nil {
			return "", err
		}

		_, err = io.Copy(out, file)

		if err != nil {
			return "", err
		}
	}

	return f.dirname + "/" + id, nil
}

func (f *fs) Delete(ctx context.Context, id string) error {
	return nil
}
