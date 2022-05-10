package miniofs

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

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
	exists, err := f.client.BucketExists(ctx, foldername)
	if err != nil {
		return "", "", err
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

func (f *fs) Delete(ctx context.Context, id string) error {
	// cause we store the whole path to a folder in our
	// database we have to extract id of a folder
	// from that path
	// TODO: would be great if we would not hav to do that
	temp := strings.Split(id, "/")
	if len(temp) < 3 {
		return errors.New("miniofs: incorrect id of bucket")
	}
	folder := temp[3]
	objectsCh := make(chan minio.ObjectInfo)
	errCh := make(chan error, 1)
	go func() {
		defer close(objectsCh)
		doneCh := make(chan struct{})
		defer close(doneCh)
		for object := range f.client.ListObjects(ctx, folder, minio.ListObjectsOptions{Prefix: "", Recursive: true}) {
			if object.Err != nil {
				errCh <- object.Err
				return
			}
			objectsCh <- object
		}
	}()
	for err := range f.client.RemoveObjects(ctx, folder, objectsCh, minio.RemoveObjectsOptions{GovernanceBypass: true}) {
		if err.Err != nil {
			return err.Err
		}
	}
	if err := f.client.RemoveBucket(ctx, folder); err != nil {
		errCh <- err
	}
	errCh <- nil
	err := <-errCh
	return err
}
