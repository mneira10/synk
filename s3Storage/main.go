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
	client     *s3.Client
	BucketName string
	Url        string
}

type S3Storage interface {
	ListObjects() *s3.ListObjectsV2Output
	// UploadFile(bucketPath string)
	// DownloadFile(bucketPath string)
}

// TODO: generalize this configuration to any S3 source
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
		client:     s3.NewFromConfig(cfg),
		BucketName: "test-synk",
		Url:        fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
	}

	log.Info("Successfully configured s3.")
	return &s3Obj
}

func (s3Obj *S3Object) ListObjects() *s3.ListObjectsV2Output {

	log.WithFields(log.Fields{"bucketName": s3Obj.BucketName}).Info("Listing objects")

	// This should work for up to 1k objects:
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.ListObjectsV2
	// TODO: get all objects here
	listObjectsOutput, err := s3Obj.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  &s3Obj.BucketName,
		MaxKeys: 1000,
	})

	if err != nil {
		log.Fatal(err)
		log.Fatal("Could not list files in bucket")
		os.Exit(1)
	}

	return listObjectsOutput

	// for _, object := range listObjectsOutput.Contents {
	// 	// obj, _ := json.MarshalIndent(object, "", "\t")
	// 	// fmt.Println(string(obj))
	// 	fmt.Printf("Name: %v\n", *object.Key)
	// }
}
