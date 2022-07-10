package s3Storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	fmt.Println("Configuring S3...")

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
		log.Printf("error: %v", err)
		return nil
	}

	fmt.Println("Returning configuration...")

	s3Obj := S3Object{
		client: s3.NewFromConfig(cfg),
	}
	return &s3Obj
}

func (s3Obj *S3Object) ListFiles() {

	// TODO: generalize this
	bucketName := "test-synk"

	// This should work for up to 1k objects:
	// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.ListObjectsV2
	listObjectsOutput, err := s3Obj.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	})

	if err != nil {
		log.Fatal(err)
	}

	for _, object := range listObjectsOutput.Contents {
		obj, _ := json.MarshalIndent(object, "", "\t")
		fmt.Println(string(obj))
	}
}
