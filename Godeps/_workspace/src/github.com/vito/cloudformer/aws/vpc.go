package aws

import (
	"net"

	"github.com/vito/cloudformer"
	"github.com/vito/cloudformer/aws/models"
)

type VPC struct {
	name string

	resources models.Resources

	model *models.VPC
}

func (vpc VPC) Network(network *net.IPNet) {
	vpc.model.CidrBlock = network.String()
}

func (vpc VPC) AttachInternetGateway(gateway cloudformer.InternetGateway) {
	vpc.resources[vpc.name+"VPCGatewayAttachment"] =
		&models.VPCGatewayAttachment{
			VpcId:             models.Ref(vpc.name + "VPC"),
			InternetGatewayId: models.Ref(gateway.(InternetGateway).name + "InternetGateway"),
		}
}

func (vpc VPC) AssociateDHCPOptions(options cloudformer.DHCPOptions) {
	vpc.resources[vpc.name+"DHCPOptions"] =
		&models.DHCPOptions{
			DomainNameServers: options.DomainNameServers,
		}

	vpc.resources[vpc.name+"VPCDHCPOptionsAssociation"] =
		&models.VPCDHCPOptionsAssociation{
			VpcId:         models.Ref(vpc.name + "VPC"),
			DhcpOptionsId: models.Ref(vpc.name + "DHCPOptions"),
		}
}

func (vpc VPC) Subnet(name string) cloudformer.Subnet {
	model := &models.Subnet{
		VpcId: models.Ref(vpc.name + "VPC"),
		Tags: []models.Tag{
			{Key: "Name", Value: name},
		},
	}

	vpc.resources[name+"Subnet"] = model

	return Subnet{
		name:      name,
		model:     model,
		resources: vpc.resources,
	}
}

func (vpc VPC) SecurityGroup(name string) cloudformer.SecurityGroup {
	model := &models.SecurityGroup{
		GroupDescription:     name,
		VpcId:                models.Ref(vpc.name + "VPC"),
		SecurityGroupIngress: []interface{}{},
		SecurityGroupEgress:  []interface{}{},
	}

	vpc.resources[name+"SecurityGroup"] = model

	return SecurityGroup{
		name:  name,
		model: model,
	}
}
