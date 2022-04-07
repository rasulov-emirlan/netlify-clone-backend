package fs

import (
	"context"
	"io/ioutil"
	"mime/multipart"
	"os"
)

type fs struct {
	dirname string
}

func NewFileSystem(dirname string) (*fs, error) {
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		if err = os.Mkdir(dirname, 0555); err != nil {
			return nil, err
		}
	}
	return &fs{}, nil
}

func (f *fs) Upload(ctx context.Context, file multipart.File, id string) (path string, err error) {
	if _, err := os.Stat(f.dirname + "/" + id); os.IsExist(err) {
		return "", err
	}
	ff, err := os.OpenFile(f.dirname+"/"+id, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return "", err
	}
	defer ff.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	if _, err = ff.Write(b); err != nil {
		return "", err
	}

	return f.dirname + "/" + id, nil
}

func (f *fs) Delete(ctx context.Context, id string) error {
	return nil
}
