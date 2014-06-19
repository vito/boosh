package aws

import (
	"github.com/vito/cloudformer"
	"github.com/vito/cloudformer/aws/models"
)

type AWSCloudFormer struct {
	Template *models.Template

	resources models.Resources
}

func New(description string) *AWSCloudFormer {
	resources := make(models.Resources)

	return &AWSCloudFormer{
		Template: &models.Template{
			AWSTemplateFormatVersion: "2010-09-09",
			Description:              description,

			Resources: resources,
			Mappings:  models.Hash{},
		},

		resources: resources,
	}
}

func (cloudFormer *AWSCloudFormer) InternetGateway(name string) cloudformer.InternetGateway {
	gateway := InternetGateway{
		name:  name,
		model: &models.InternetGateway{},
	}

	cloudFormer.resources[name+"InternetGateway"] = gateway.model

	return gateway
}

func (cloudFormer *AWSCloudFormer) VPC(name string) cloudformer.VPC {
	vpc := VPC{
		name:      name,
		model:     &models.VPC{},
		resources: cloudFormer.resources,
	}

	cloudFormer.resources[name+"VPC"] = vpc.model

	return vpc
}

func (cloudFormer *AWSCloudFormer) Bucket(name string) cloudformer.Bucket {
	s3Bucket := Bucket{
		name:  name,
		model: &models.Bucket{},
	}

	cloudFormer.resources[name+"Bucket"] = s3Bucket.model

	return s3Bucket
}

func (cloudFormer *AWSCloudFormer) ElasticIP(name string) cloudformer.ElasticIP {
	ip := ElasticIP{
		name:      name,
		model:     &models.EIP{},
		resources: cloudFormer.resources,
	}

	cloudFormer.resources[name+"EIP"] = ip.model

	return ip
}

func (cloudFormer *AWSCloudFormer) LoadBalancer(name string) cloudformer.LoadBalancer {
	model := &models.LoadBalancer{
		Listeners:      []interface{}{},
		Subnets:        []interface{}{},
		SecurityGroups: []interface{}{},
	}

	cloudFormer.resources[name+"LoadBalancer"] = model

	return LoadBalancer{
		name:      name,
		model:     model,
		resources: cloudFormer.resources,
	}
}
