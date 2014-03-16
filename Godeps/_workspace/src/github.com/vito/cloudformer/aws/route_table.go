package aws

import (
	"github.com/vito/cloudformer"
	"github.com/vito/cloudformer/aws/models"
)

type RouteTable struct {
	name      string
	model     *models.RouteTable
	resources models.Resources
}

func (routeTable RouteTable) InternetGateway(gateway cloudformer.InternetGateway) {
	route := &models.Route{
		DestinationCidrBlock: "0.0.0.0/0",
		RouteTableId:         models.Ref(routeTable.name + "RouteTable"),
		GatewayId:            models.Ref(gateway.(InternetGateway).name + "InternetGateway"),
		Depends:              gateway.(InternetGateway).name + "InternetGateway",
	}

	routeTable.resources[routeTable.name+"Route"] = route
}

func (routeTable RouteTable) Instance(gateway cloudformer.Instance) {
	route := &models.Route{
		DestinationCidrBlock: "0.0.0.0/0",
		RouteTableId:         models.Ref(routeTable.name + "RouteTable"),
		InstanceId:           models.Ref(gateway.(Instance).name + "Instance"),
		Depends:              gateway.(Instance).name + "Instance",
	}

	routeTable.resources[routeTable.name+"Route"] = route
}
