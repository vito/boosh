package aws

import (
	"fmt"

	"github.com/vito/cloudformer"
	"github.com/vito/cloudformer/aws/models"
)

type LoadBalancer struct {
	name      string
	model     *models.LoadBalancer
	resources models.Resources
}

func (balancer LoadBalancer) Listener(
	protocol cloudformer.ProtocolType,
	port uint16,
	destinationProtocol cloudformer.ProtocolType,
	destinationPort uint16,
	sslCertificateId string,
) {
	listeners := balancer.model.Listeners.([]interface{})

	listeners = append(listeners, models.LoadBalancerListener{
		Protocol:         string(protocol),
		LoadBalancerPort: fmt.Sprintf("%d", port),
		InstanceProtocol: string(destinationProtocol),
		InstancePort:     fmt.Sprintf("%d", destinationPort),
		SSLCertificateId: sslCertificateId,
	})

	balancer.model.Listeners = listeners
}

func (balancer LoadBalancer) HealthCheck(check cloudformer.HealthCheck) {
	balancer.model.HealthCheck = models.LoadBalancerHealthCheck{
		Target:             fmt.Sprintf("%s:%d%s", check.Protocol, check.Port, check.Path),
		Interval:           fmt.Sprintf("%d", int(check.Interval.Seconds())),
		Timeout:            fmt.Sprintf("%d", int(check.Timeout.Seconds())),
		HealthyThreshold:   fmt.Sprintf("%d", check.HealthyThreshold),
		UnhealthyThreshold: fmt.Sprintf("%d", check.UnhealthyThreshold),
	}
}

func (balancer LoadBalancer) Scheme(scheme string) {
	balancer.model.Scheme = scheme
}

func (balancer LoadBalancer) Subnet(subnet cloudformer.Subnet) {
	subnets := balancer.model.Subnets.([]interface{})

	subnets = append(
		subnets,
		models.Ref(subnet.(Subnet).name+"Subnet"),
	)

	balancer.model.Subnets = subnets
}

func (balancer LoadBalancer) SecurityGroup(group cloudformer.SecurityGroup) {
	securityGroups := balancer.model.SecurityGroups.([]interface{})

	securityGroups = append(
		securityGroups,
		models.Ref(group.(SecurityGroup).name+"SecurityGroup"),
	)

	balancer.model.SecurityGroups = securityGroups
}

func (balancer LoadBalancer) RecordSet(name, zone string) {
	balancer.resources[balancer.name+"RecordSet"] =
		&models.RecordSetGroup{
			HostedZoneName: zone + ".",
			RecordSets: []models.RecordSet{
				{
					Name: name + "." + zone + ".",
					Type: "A",
					AliasTarget: models.RecordSetAliasTarget{
						HostedZoneId: models.Hash{
							"Fn::GetAtt": []string{
								balancer.name + "LoadBalancer",
								"CanonicalHostedZoneNameID",
							},
						},
						DNSName: models.Hash{
							"Fn::GetAtt": []string{
								balancer.name + "LoadBalancer",
								"CanonicalHostedZoneName",
							},
						},
					},
				},
			},
		}
}
