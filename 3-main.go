package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/fraenkel/candiedyaml"
)

// amzn-ami-vpc-nat-pv-2013.09.0.x86_64-ebs
var NAT_AMIS = map[string]string{
	"us-east-1":      "ami-ad227cc4",
	"us-west-1":      "ami-d69aad93",
	"us-west-2":      "ami-f032acc0",
	"eu-west-1":      "ami-f3e30084",
	"ap-southeast-1": "ami-f22772a0",
	"ap-southeast-2": "ami-3bae3201",
	"ap-northeast-1": "ami-cd43d9cc",
	"sa-east-1":      "ami-d78325ca",
}

func main() {
	var spec DeploymentSpec

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	candiedyaml.NewDecoder(file).Decode(&spec)

	resources := make(map[string]Typer)

	resources["VPC"] = Vpc{
		CidrBlock: spec.VPC.CIDR,
	}

	resources["VPCInternetGateway"] = VPCGatewayAttachment{
		InternetGatewayId: ref(spec.VPC.InternetGateway + "InternetGateway"),
		VpcId:             ref("VPC"),
	}

	if len(spec.DNS) > 0 {
		resources["DHCPOptions"] = DHCPOptions{
			DomainNameServers: spec.DNS,
		}

		resources["VPCDHCPOptionsAssociation"] = VPCDHCPOptionsAssociation{
			VpcId:         ref("VPC"),
			DhcpOptionsId: ref("DHCPOptions"),
		}
	}

	for _, x := range spec.InternetGateways {
		resources[x.Name+"InternetGateway"] = InternetGateway{}
	}

	for _, x := range spec.Subnets {
		resources[x.Name+"Subnet"] = Subnet{
			AvailabilityZone: x.AvailabilityZone,
			CidrBlock:        x.CIDR,
			VpcId:            ref("VPC"),
			Tags: []interface{}{
				Tag{Key: "Name", Value: x.Name},
			},
		}

		if x.RouteTable != nil {
			resources[x.Name+"SubnetRouteTable"] = RouteTable{
				VpcId: ref("VPC"),
			}

			if x.RouteTable.Instance != nil {
				resources[x.Name+"SubnetRoute"] = Route{
					DestinationCidrBlock: "0.0.0.0/0",
					RouteTableId:         ref(x.Name + "SubnetRouteTable"),
					InstanceId:           ref(*x.RouteTable.Instance + "NATInstance"),
				}
			} else if x.RouteTable.InternetGateway != nil {
				resources[x.Name+"SubnetRoute"] = Route{
					DestinationCidrBlock: "0.0.0.0/0",
					RouteTableId:         ref(x.Name + "SubnetRouteTable"),
					GatewayId:            ref(*x.RouteTable.InternetGateway + "InternetGateway"),
				}
			}

			resources[x.Name+"SubnetRouteTableAssociation"] = SubnetRouteTableAssociation{
				RouteTableId: ref(x.Name + "SubnetRouteTable"),
				SubnetId:     ref(x.Name + "Subnet"),
			}
		}

		if x.NAT != nil {
			resources[x.NAT.Name+"NATInstance"] = Instance{
				AvailabilityZone: x.AvailabilityZone,
				InstanceType:     x.NAT.InstanceType,
				PrivateIpAddress: x.NAT.IP,
				KeyName:          x.NAT.KeyPairName,
				SubnetId:         ref(x.Name + "Subnet"),
				ImageId: Hash{
					"Fn::FindInMap": []interface{}{
						"AWSNATAMI",
						ref("AWS::Region"),
						"AMI",
					},
				},
				SecurityGroupIds: []interface{}{
					ref(x.NAT.SecurityGroup + "SecurityGroup"),
				},
				Tags: []interface{}{
					Tag{Key: "Name", Value: x.NAT.Name},
				},
			}
		}
	}

	for _, x := range spec.SecurityGroups {
		ingress := []interface{}{}
		egress := []interface{}{}

		for _, i := range x.Ingress {
			fromPort, toPort := parsePortRange(i.Ports)

			ingress = append(ingress, SecurityGroupIngress{
				CidrIp:     i.CIDR,
				IpProtocol: i.Protocol,
				FromPort:   fromPort,
				ToPort:     toPort,
			})
		}

		for _, e := range x.Egress {
			fromPort, toPort := parsePortRange(e.Ports)

			ingress = append(ingress, SecurityGroupEgress{
				CidrIp:     e.CIDR,
				IpProtocol: e.Protocol,
				FromPort:   fromPort,
				ToPort:     toPort,
			})
		}

		resources[x.Name+"SecurityGroup"] = SecurityGroup{
			GroupDescription:     x.Name,
			VpcId:                ref("VPC"),
			SecurityGroupIngress: ingress,
			SecurityGroupEgress:  egress,
		}
	}

	for _, x := range spec.LoadBalancers {
		subnets := []interface{}{}
		for _, name := range x.Subnets {
			subnets = append(subnets, ref(name+"Subnet"))
		}

		listeners := []interface{}{}
		for _, listener := range x.Listeners {
			destinationPort := listener.Port
			if listener.DestinationPort != nil {
				destinationPort = *listener.DestinationPort
			}

			destinationProtocol := listener.Protocol
			if listener.DestinationProtocol != nil {
				destinationProtocol = *listener.DestinationProtocol
			}

			listeners = append(listeners, LoadBalancerListener{
				LoadBalancerPort: listener.Port,
				Protocol:         listener.Protocol,
				InstancePort:     destinationPort,
				InstanceProtocol: destinationProtocol,
			})
		}

		resources[x.Name+"LoadBalancer"] = LoadBalancer{
			Subnets:   subnets,
			Listeners: listeners,
			HealthCheck: LoadBalancerHealthCheck{
				Target:             x.HealthCheck.Target.Type + ":" + x.HealthCheck.Target.Port,
				Timeout:            x.HealthCheck.Timeout,
				Interval:           x.HealthCheck.Interval,
				HealthyThreshold:   x.HealthCheck.HealthyThreshold,
				UnhealthyThreshold: x.HealthCheck.UnhealthyThreshold,
			},
		}
	}

	template := &Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Description:              "lol",

		Resources: resources,

		Mappings: Hash{
			"AWSNATAMI": Hash{
				"us-east-1": Hash{
					"AMI": NAT_AMIS["us-east-1"],
				},
				"us-west-1": Hash{
					"AMI": NAT_AMIS["us-west-1"],
				},
				"us-west-2": Hash{
					"AMI": NAT_AMIS["us-west-2"],
				},
				"eu-west-1": Hash{
					"AMI": NAT_AMIS["eu-west-1"],
				},
				"ap-southeast-1": Hash{
					"AMI": NAT_AMIS["ap-southeast-1"],
				},
				"ap-southeast-2": Hash{
					"AMI": NAT_AMIS["ap-southeast-2"],
				},
				"ap-northeast-1": Hash{
					"AMI": NAT_AMIS["ap-northeast-1"],
				},
				"sa-east-1": Hash{
					"AMI": NAT_AMIS["sa-east-1"],
				},
			},
		},
	}

	json.NewEncoder(os.Stdout).Encode(template)
}

func parsePortRange(ports string) (string, string) {
	segments := strings.Split(ports, "-")

	fromPort := ""
	toPort := ""

	if len(segments) == 1 {
		fromPort = segments[0]
		toPort = fromPort
	} else if len(segments) == 2 {
		fromPort = segments[0]
		toPort = segments[1]
	}

	fromPort = strings.Trim(fromPort, " ")
	toPort = strings.Trim(toPort, " ")

	return fromPort, toPort
}
