package oss

import (
	"bytes"
	"io"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func PutObject(endpoint, accessKeyID, accessKeySecret, bucketName, objectKey string, val []byte) error {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}
	err = bucket.PutObject(objectKey, bytes.NewReader(val))
	if err != nil {
		return err
	}
	return nil
}

func GetObjectSignedUrl(endpoint, accessKeyID, accessKeySecret, bucketName, objectKey string, expiredInSec int64, https bool) (string, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", err
	}
	signedURL, err := bucket.SignURL(objectKey, oss.HTTPGet, expiredInSec)
	if err != nil {
		return "", err
	}
	if !https {
		return signedURL, nil
	}
	return strings.Replace(signedURL, "http:", "https:", 1), nil
}

func GetObject(endpoint, accessKeyID, accessKeySecret, bucketName, objectKey string) ([]byte, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	object, err := bucket.GetObject(objectKey)
	if err != nil {
		return nil, err
	}
	defer object.Close()
	data, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}
	return data, nil
}
