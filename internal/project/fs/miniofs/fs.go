package miniofs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type fs struct {
	client *minio.Client
	url    string
}

func NewFileSystem(url, accessKeyID, secretAccessKey string, useSSL bool) (*fs, error) {
	c, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &fs{client: c, url: "http://" + url}, nil
}

func (f *fs) Upload(ctx context.Context, files []*multipart.FileHeader, id string) (path string, err error) {
	err = f.client.MakeBucket(ctx, id, minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: false})
	if err != nil {
		exists, errBucketExists := f.client.BucketExists(ctx, id)
		if errBucketExists == nil && exists {
			return "", errors.New("miniofs: trying to create a bucket that already exists")
		}
		return "", err
	}
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{"Action": ["s3:GetObject"],
		"Effect": "Allow",
		"Principal": {"AWS": ["*"]},
		"Resource": ["arn:aws:s3:::%s/*"],
		"Sid": ""}]}`, id)
	if err := f.client.SetBucketPolicy(ctx, id, policy); err != nil {
		return path, err
	}
	for _, v := range files {
		file, err := v.Open()
		if err != nil {
			return "", err
		}
		defer file.Close()
		log.Println("trying to upload files")
		if _, err = f.client.PutObject(
			ctx,
			id,
			v.Filename,
			file,
			v.Size,
			minio.PutObjectOptions{ContentType: v.Header.Get("Content-Type")},
		); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s/%s/", f.url, id), nil
}

// Be careful with this function. It deletes all the files inside of a bucket
// and replaces them with new ones
func (f *fs) Replace(ctx context.Context, files []*multipart.FileHeader, foldername string) error {
	exists, err := f.client.BucketExists(ctx, foldername)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("miniofs: we don't have such bucket")
	}
	objectsCh := make(chan minio.ObjectInfo)
	errCh := make(chan error)
	go func() {
		defer close(objectsCh)
		doneCh := make(chan struct{})
		defer close(doneCh)
		for object := range f.client.ListObjects(ctx, "mytestbucket", minio.ListObjectsOptions{Prefix: "", Recursive: true}) {
			if object.Err != nil {
				errCh <- object.Err
				return
			}
			objectsCh <- object
		}
	}()
	go func() {
		errorCh := f.client.RemoveObjects(ctx, foldername, objectsCh, minio.RemoveObjectsOptions{})
		for e := range errorCh {
			errCh <- e.Err
			return
		}
	}()
	select {
	case <-errCh:
		return <-errCh
	}
}

func (f *fs) Delete(ctx context.Context, id string) error {
	panic("not implemented")
}
