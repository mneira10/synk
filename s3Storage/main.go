package s3Storage

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/mneira10/synk/logger"
)

type S3Object struct {
	client *s3.Client
}

type S3Storage interface {
	ListItems()
	UploadFile(bucketPath string)
	DownloadFile(bucketPath string)
}

func ConfigS3() *S3Object {
	log.Info("Configuring S3...")

	// temp hack while I implement a config file somewhere
	accountId := os.Getenv("CLOUDFLARE_R2_ACCOUNT_ID")
	accessKeyId := os.Getenv("CLOUDFLARE_R2_ACCESS_KEY")
	accessKeySecret := os.Getenv("CLOUDFLARE_R2_SECRET_KEY")

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
	)

	if err != nil {
		log.Error(fmt.Printf("error: %v", err))
		return nil
	}

	s3Obj := S3Object{
		client: s3.NewFromConfig(cfg),
	}

	log.Info("Successfully configured s3.")
	return &s3Obj
}

func (s3Obj *S3Object) ListFiles() {

	// TODO: generalize this
	bucketName := "test-synk"

	log.WithFields(log.Fields{"bucketName": bucketName}).Info("Listing objects")

	// This should work for up to 1k objects:
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.ListObjectsV2
	// TODO: get all objects here
	listObjectsOutput, err := s3Obj.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  &bucketName,
		MaxKeys: 1000,
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, object := range listObjectsOutput.Contents {
		// obj, _ := json.MarshalIndent(object, "", "\t")
		// fmt.Println(string(obj))
		fmt.Printf("Name: %v\n", *object.Key)
	}
}
