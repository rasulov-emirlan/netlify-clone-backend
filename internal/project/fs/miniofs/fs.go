package miniofs

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type fs struct {
	client     *minio.Client
	url        string
	baseBucket string
}

func NewFileSystem(url, accessKeyID, secretAccessKey, baseBucket string, useSSL bool) (*fs, error) {
	c, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &fs{client: c, url: "http://" + url, baseBucket: baseBucket}, nil
}

func (f *fs) Upload(ctx context.Context, files []*multipart.FileHeader, foldername string, version int) (string, string, error) {
	foldername = fmt.Sprintf("%s-%s", f.baseBucket, foldername)
	exists, errBucketExists := f.client.BucketExists(ctx, foldername)
	if errBucketExists == nil && !exists {
		return "", "", errBucketExists
	}
	if !exists {
		if err := f.client.MakeBucket(
			ctx,
			foldername,
			minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: false},
		); err != nil {
			return "", "", err
		}
		// i dont' know how code bellow works
		// so please dont touch it
		policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{"Action": ["s3:GetObject"],
		"Effect": "Allow",
		"Principal": {"AWS": ["*"]},
		"Resource": ["arn:aws:s3:::%s/*"],
		"Sid": ""}]}`, foldername)
		if err := f.client.SetBucketPolicy(ctx, foldername, policy); err != nil {
			return "", "", err
		}
	}

	for _, v := range files {
		file, err := v.Open()
		if err != nil {
			return "", "", err
		}
		defer file.Close()
		filename := v.Filename
		if ext := filepath.Ext(filename); ext == ".js" ||
			ext == ".html" || ext == ".css" {
			filename = fmt.Sprintf("%d/%s", version, filename)
		}
		if _, err = f.client.PutObject(
			ctx,
			foldername,
			filename,
			file,
			v.Size,
			minio.PutObjectOptions{ContentType: v.Header.Get("Content-Type")},
		); err != nil {
			return "", "", err
		}
	}
	return fmt.Sprintf("%s/%s/%d/", f.url, foldername, version), fmt.Sprintf("%s/%s/", f.url, foldername), nil
}

// this function is deprecated and it is not finished
// func (f *fs) Replace(ctx context.Context, files []*multipart.FileHeader, foldername string) error {
// 	foldername = fmt.Sprintf("%s-%s", f.baseBucket, foldername)
// 	exists, err := f.client.BucketExists(ctx, foldername)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return errors.New("miniofs: we don't have such bucket")
// 	}
// 	objectsCh := make(chan minio.ObjectInfo)
// 	errCh := make(chan error)
// 	go func() {
// 		defer close(objectsCh)
// 		doneCh := make(chan struct{})
// 		defer close(doneCh)
// 		for object := range f.client.ListObjects(ctx, "mytestbucket", minio.ListObjectsOptions{Prefix: "", Recursive: true}) {
// 			if object.Err != nil {
// 				errCh <- object.Err
// 				return
// 			}
// 			objectsCh <- object
// 		}
// 	}()
// 	go func() {
// 		errorCh := f.client.RemoveObjects(ctx, foldername, objectsCh, minio.RemoveObjectsOptions{})
// 		for e := range errorCh {
// 			errCh <- e.Err
// 			return
// 		}
// 	}()
// 	select {
// 	case <-errCh:
// 		return <-errCh
// 	}
// }

func (f *fs) Delete(ctx context.Context, id string) error {
	panic("not implemented")
}
