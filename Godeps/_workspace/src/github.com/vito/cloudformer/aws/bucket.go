package aws

import "github.com/vito/cloudformer/aws/models"

type Bucket struct {
	name string

	model *models.Bucket
}

func (bucket Bucket) Name(name string) {
	bucket.model.BucketName = name
}
