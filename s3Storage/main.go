package s3Storage

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	log "github.com/mneira10/synk/logger"
)

type S3Object struct {
	client     *s3.Client
	config     aws.Config
	BucketName string
	Url        string
}

type S3Storage interface {
	ListObjects() []types.Object
	UploadFile(localFilePath string, bucketPath string) error
	DeleteFile(bucketFilePath string) error
	// DownloadFile(bucketPath string)
}

// TODO: generalize this configuration to any S3 source
func ConfigS3(storageConfig *R2ConfigData) *S3Object {
	log.Info("Configuring S3...")

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", storageConfig.AccountId),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				storageConfig.AccessKeyId,
				storageConfig.AccessKeySecret,
				"")),
	)

	if err != nil {
		log.Error(fmt.Printf("error: %v", err))
		return nil
	}

	s3Obj := S3Object{
		client:     s3.NewFromConfig(cfg),
		config:     cfg,
		BucketName: storageConfig.BucketName,
		Url:        fmt.Sprintf("https://%s.r2.cloudflarestorage.com", storageConfig.AccountId),
	}

	log.Info("Successfully configured s3.")
	return &s3Obj
}

func (s3Obj *S3Object) ListObjects() []types.Object {

	log.WithFields(log.Fields{"bucketName": s3Obj.BucketName}).Info("Listing objects")

	// This should work for up to 1k objects:
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.ListObjectsV2
	// TODO: get all objects here
	listObjectsOutput, err := s3Obj.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  &s3Obj.BucketName,
		MaxKeys: 1000,
	})

	if err != nil {
		fmt.Printf("Could not list files in %v. Double check your configuration!\n", s3Obj.BucketName)
		log.Error("Could not list files in bucket")
		log.Fatal(err)
		os.Exit(1)
	}

	allContents := listObjectsOutput.Contents

	isTruncated := listObjectsOutput.IsTruncated
	startAfter := listObjectsOutput.StartAfter

	for isTruncated {
		newListObjectsOutput, err := listObjects(&s3Obj.BucketName, startAfter, s3Obj)

		if err != nil {
			fmt.Printf("Could not list files in %v. Double check your configuration!\n", s3Obj.BucketName)
			log.Error("Could not list files in bucket")
			log.Fatal(err)
			os.Exit(1)
		}

		allContents = append(allContents, newListObjectsOutput.Contents...)

		isTruncated = newListObjectsOutput.IsTruncated
		startAfter = newListObjectsOutput.StartAfter

	}

	return allContents
}

func listObjects(bucketName *string, startAfter *string, s3Obj *S3Object) (*s3.ListObjectsV2Output, error) {
	log.WithFields(log.Fields{"startAfter": startAfter}).Info("Listing again after")
	listObjectsOutput, err := s3Obj.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  bucketName,
		MaxKeys: 1000,
	})

	if err != nil {
		fmt.Printf("Could not list files in %v. Double check your configuration!\n", s3Obj.BucketName)
		log.Error("Could not list files in bucket")
		log.Fatal(err)
		os.Exit(1)
	}

	return listObjectsOutput, err
}

func (s3Obj *S3Object) UploadFile(localFilePath string, bucketPath string) error {

	upFile, err := os.Open(localFilePath)

	if err != nil {
		return fmt.Errorf("could not open local filepath [%v]: %+v", localFilePath, err)
	}

	defer upFile.Close()

	// Get the file info
	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	objectData := &s3.PutObjectInput{
		Bucket:        aws.String(s3Obj.BucketName),
		Key:           aws.String(bucketPath),
		Body:          bytes.NewReader(fileBuffer),
		ContentLength: fileSize,
		ContentType:   aws.String(http.DetectContentType(fileBuffer)),
		// TODO: look into this
		// ACL:                  "private",
		// ContentDisposition:   aws.String("attachment"),
		// ServerSideEncryption: "AES256",
	}

	_, err = s3Obj.client.PutObject(context.TODO(), objectData)
	return err
}

func (s3Obj *S3Object) DeleteFile(bucketFilePath string) error {

	deleteObjData := &s3.DeleteObjectInput{
		Bucket: &s3Obj.BucketName,
		Key:    &bucketFilePath,
	}

	_, err := s3Obj.client.DeleteObject(context.TODO(), deleteObjData)
	return err
}
