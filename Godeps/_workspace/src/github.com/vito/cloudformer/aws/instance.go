package aws

import (
	"github.com/vito/cloudformer"
	"net"

	"github.com/vito/cloudformer/aws/models"
)

type Instance struct {
	name  string
	model *models.Instance
}

func (instance Instance) Type(name string) {
	instance.model.InstanceType = name
}

func (instance Instance) Image(id string) {
	instance.model.ImageId = id
}

func (instance Instance) PrivateIP(ip net.IP) {
	instance.model.PrivateIpAddress = ip.String()
}

func (instance Instance) KeyPair(name string) {
	instance.model.KeyName = name
}

func (instance Instance) SecurityGroup(securityGroup cloudformer.SecurityGroup) {
	groups := instance.model.SecurityGroupIds.([]interface{})

	groups = append(
		groups,
		models.Ref(securityGroup.(SecurityGroup).name+"SecurityGroup"),
	)

	instance.model.SecurityGroupIds = groups
}

func (instance Instance) SourceDestCheck(val bool) {
	instance.model.SourceDestCheck = val
}
