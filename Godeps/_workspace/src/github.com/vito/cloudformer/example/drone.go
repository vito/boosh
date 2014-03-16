package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/vito/cloudformer"
	"github.com/vito/cloudformer/aws"
)

func Form(f cloudformer.CloudFormer) {
	zone1 := "us-east-1a"

	vpc := f.VPC("Drone")
	vpc.Network(cloudformer.CIDR("10.10.0.0/16"))

	vpcGateway := f.InternetGateway("Drone")

	vpc.AttachInternetGateway(vpcGateway)

	openSecurityGroup := vpc.SecurityGroup("Open")
	boshSecurityGroup := vpc.SecurityGroup("BOSH")
	internalSecurityGroup := vpc.SecurityGroup("Internal")
	webSecurityGroup := vpc.SecurityGroup("Web")

	for _, group := range []cloudformer.SecurityGroup{
		openSecurityGroup,
		boshSecurityGroup,
		internalSecurityGroup,
	} {
		group.Ingress(cloudformer.TCP, cloudformer.CIDR("0.0.0.0/0"), 0, 65535)
		group.Ingress(cloudformer.UDP, cloudformer.CIDR("0.0.0.0/0"), 0, 65535)
	}

	webSecurityGroup.Ingress(cloudformer.TCP, cloudformer.CIDR("0.0.0.0/0"), 80, 80)
	webSecurityGroup.Ingress(cloudformer.TCP, cloudformer.CIDR("0.0.0.0/0"), 8080, 8080)

	boshSubnet := vpc.Subnet("BOSH")
	boshSubnet.Network(cloudformer.CIDR("10.10.0.0/24"))
	boshSubnet.AvailabilityZone(zone1)
	boshSubnet.RouteTable().InternetGateway(vpcGateway)

	droneELBSubnet := vpc.Subnet("DroneELB")
	droneELBSubnet.Network(cloudformer.CIDR("10.10.2.0/24"))
	droneELBSubnet.AvailabilityZone(zone1)
	droneELBSubnet.RouteTable().InternetGateway(vpcGateway)

	droneSubnet := vpc.Subnet("Drone")
	droneSubnet.Network(cloudformer.CIDR("10.10.16.0/20"))
	droneSubnet.AvailabilityZone(zone1)

	boshNAT := boshSubnet.Instance("NAT")
	boshNAT.Type("m1.small")
	boshNAT.Image("ami-something")
	boshNAT.PrivateIP(cloudformer.IP("10.10.0.10"))
	boshNAT.KeyPair("bosh")
	boshNAT.SecurityGroup(openSecurityGroup)

	droneSubnet.RouteTable().Instance(boshNAT)

	balancer := f.LoadBalancer("Drone")
	balancer.Listener(cloudformer.TCP, 80, cloudformer.TCP, 80)
	balancer.Listener(cloudformer.TCP, 8080, cloudformer.TCP, 8080)
	balancer.HealthCheck(cloudformer.HealthCheck{
		Protocol:           cloudformer.TCP,
		Port:               80,
		Timeout:            5 * time.Second,
		Interval:           30 * time.Second,
		HealthyThreshold:   10,
		UnhealthyThreshold: 2,
	})
	balancer.Subnet(droneELBSubnet)
	balancer.SecurityGroup(webSecurityGroup)
}

func main() {
	aws := aws.New("test template")

	Form(aws)

	json.NewEncoder(os.Stdout).Encode(aws.Template)
}
