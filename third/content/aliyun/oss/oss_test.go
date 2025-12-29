package oss

import (
	"os"
	"testing"
)

func TestAll(t *testing.T) {
	TestPutObject(t)
	TestGetObjectSignedUrl(t)
	TestGetObject(t)
}

func TestPutObject(t *testing.T) {
	bucketName := os.Getenv("BucketName")
	err := PutObject(
		os.Getenv("Endpoint"),
		os.Getenv("AccessKeyID"),
		os.Getenv("AccessKeySecret"),
		bucketName,
		"file.txt",
		[]byte(bucketName),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("OK")
}

func TestGetObjectSignedUrl(t *testing.T) {
	bucketName := os.Getenv("BucketName")
	url, err := GetObjectSignedUrl(
		os.Getenv("Endpoint"),
		os.Getenv("AccessKeyID"),
		os.Getenv("AccessKeySecret"),
		bucketName,
		"file.txt",
		3600,
		true,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}

func TestGetObject(t *testing.T) {
	bucketName := os.Getenv("BucketName")
	result, err := GetObject(
		os.Getenv("Endpoint"),
		os.Getenv("AccessKeyID"),
		os.Getenv("AccessKeySecret"),
		bucketName,
		"file.txt",
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(result))
}
