package aws

import (
	"fmt"
	"net"

	"github.com/vito/cloudformer"
	"github.com/vito/cloudformer/aws/models"
)

type SecurityGroup struct {
	name  string
	model *models.SecurityGroup
}

func (securityGroup SecurityGroup) Ingress(
	protocol cloudformer.ProtocolType,
	network *net.IPNet,
	fromPort uint16,
	toPort uint16,
) {
	ingress := securityGroup.model.SecurityGroupIngress.([]interface{})

	ingress = append(ingress,
		&models.SecurityGroupIngress{
			CidrIp:     network.String(),
			IpProtocol: protocol,
			FromPort:   fmt.Sprintf("%d", fromPort),
			ToPort:     fmt.Sprintf("%d", toPort),
		})

	securityGroup.model.SecurityGroupIngress = ingress
}

func (securityGroup SecurityGroup) Egress(
	protocol cloudformer.ProtocolType,
	network *net.IPNet,
	fromPort uint16,
	toPort uint16,
) {
	e := securityGroup.model.SecurityGroupEgress.([]interface{})

	e = append(e,
		&models.SecurityGroupEgress{
			CidrIp:     network.String(),
			IpProtocol: protocol,
			FromPort:   fmt.Sprintf("%d", fromPort),
			ToPort:     fmt.Sprintf("%d", toPort),
		})

	securityGroup.model.SecurityGroupEgress = e
}
