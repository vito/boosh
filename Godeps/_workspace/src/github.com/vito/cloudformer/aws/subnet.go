package aws

import (
	"github.com/vito/cloudformer"
	"net"

	"github.com/vito/cloudformer/aws/models"
)

type Subnet struct {
	name      string
	model     *models.Subnet
	resources models.Resources
}

func (subnet Subnet) Network(network *net.IPNet) {
	subnet.model.CidrBlock = network.String()
}

func (subnet Subnet) AvailabilityZone(name string) {
	subnet.model.AvailabilityZone = name
}

func (subnet Subnet) Instance(name string) cloudformer.Instance {
	if subnet.model.AvailabilityZone == "" {
		panic("must set availability zone on subnet before creating an instance")
	}

	model := &models.Instance{
		AvailabilityZone: subnet.model.AvailabilityZone,
		SubnetId:         models.Ref(subnet.name + "Subnet"),
		SecurityGroupIds: []interface{}{},
		Tags: []models.Tag{
			{Key: "Name", Value: name},
		},
	}

	subnet.resources[name+"Instance"] = model

	return Instance{
		name:  name,
		model: model,
	}
}

func (subnet Subnet) RouteTable() cloudformer.RouteTable {
	model := &models.RouteTable{
		VpcId: subnet.model.VpcId,
	}

	subnet.resources[subnet.name+"RouteTable"] = model

	subnet.resources[subnet.name+"SubnetRouteTableAssociation"] =
		&models.SubnetRouteTableAssociation{
			SubnetId:     models.Ref(subnet.name + "Subnet"),
			RouteTableId: models.Ref(subnet.name + "RouteTable"),
		}

	return RouteTable{
		name:      subnet.name,
		model:     model,
		resources: subnet.resources,
	}
}
