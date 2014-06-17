package aws

import "github.com/vito/cloudformer/aws/models"

type S3Bucket struct {
	name string

	model *models.S3Bucket
}

func (s3Bucket S3Bucket) Name(name string) {
	s3Bucket.model.BucketName = name
}
