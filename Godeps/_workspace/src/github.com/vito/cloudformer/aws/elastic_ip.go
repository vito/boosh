package aws

import (
	"github.com/vito/cloudformer"
	"github.com/vito/cloudformer/aws/models"
)

type ElasticIP struct {
	name      string
	model     *models.EIP
	resources models.Resources
}

func (elasticIP ElasticIP) Domain(name string) {
	elasticIP.model.Domain = name
}

func (elasticIP ElasticIP) AttachTo(instance cloudformer.Instance) {
	elasticIP.model.InstanceId =
		models.Ref(instance.(Instance).name + "Instance")
}
