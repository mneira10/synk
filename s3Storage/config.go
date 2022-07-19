package s3Storage

// Currently only supports Cloudflare's R2
type R2ConfigData struct {
	Type            string `validate:"required,oneof='R2'"`
	BucketName      string `validate:"required"`
	Url             string `validate:"required"`
	AccountId       string `validate:"required"`
	AccessKeyId     string `validate:"required"`
	AccessKeySecret string `validate:"required"`
}
