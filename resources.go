package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/dynport/gocloud/aws/cloudformation"
	"github.com/dynport/gocloud/aws/ec2"
	"github.com/dynport/gocloud/aws/elb"
	"github.com/fraenkel/candiedyaml"
)

func resources(name string) {
	cf := cloudformation.NewFromEnv()
	ec2Client := ec2.NewFromEnv()
	elbClient := elb.NewFromEnv()

	resources, err := cf.DescribeStackResources(
		cloudformation.DescribeStackResourcesParameters{
			StackName: name,
		},
	)
	if err != nil {
		fatal(err)
	}

	stub := make(map[string]map[string]interface{})

	for _, resource := range resources.DescribeStackResourcesResult.StackResources {
		typeSegments := strings.Split(resource.ResourceType, "::")
		typeBase := typeSegments[len(typeSegments)-1]

		byType, found := stub[typeBase]
		if !found {
			byType = make(map[string]interface{})
			stub[typeBase] = byType
		}

		cleanName := strings.Replace(resource.LogicalResourceId, typeBase, "", 1)
		if cleanName == "" {
			cleanName = resource.LogicalResourceId
		}

		byType[cleanName] = resource.PhysicalResourceId
	}

	err = grabSecurityGroupNames(stub, ec2Client)
	if err != nil {
		fatal(err)
	}

	err = grabSubnetInfo(stub, ec2Client)
	if err != nil {
		fatal(err)
	}

	err = grabLoadBalancerDNSNames(stub, elbClient)
	if err != nil {
		fatal(err)
	}

	err = candiedyaml.NewEncoder(os.Stdout).Encode(map[string]interface{}{
		"Region":          ec2Client.Client.Region,
		"AccessKeyID":     ec2Client.Client.Key,
		"SecretAccessKey": ec2Client.Client.Secret,
		"Resources":       stub,
	})
	if err != nil {
		fatal(err)
	}
}

func grabSecurityGroupNames(stub map[string]map[string]interface{}, ec2Client *ec2.Client) error {
	names := make(map[string]interface{})
	stub["SecurityGroupName"] = names

	securityGroupIds := []string{}
	nameForId := make(map[string]string)
	for name, id := range stub["SecurityGroup"] {
		securityGroupIds = append(securityGroupIds, id.(string))
		nameForId[id.(string)] = name
	}

	groups, err := ec2Client.DescribeSecurityGroups(
		&ec2.DescribeSecurityGroupsParameters{
			GroupIds: securityGroupIds,
		},
	)
	if err != nil {
		return err
	}

	for _, group := range groups {
		name := nameForId[group.GroupId]
		names[name] = group.GroupName
	}

	return nil
}

func grabLoadBalancerDNSNames(stub map[string]map[string]interface{}, elbClient *elb.Client) error {
	lbs, err := elbClient.DescribeLoadBalancers()
	if err != nil {
		return err
	}

	allNames := map[string]string{}
	for _, lb := range lbs {
		allNames[lb.LoadBalancerName] = lb.DNSName
	}

	stackNames := map[string]interface{}{}
	for name, id := range stub["LoadBalancer"] {
		stackNames[name] = allNames[id.(string)]
	}

	stub["LoadBalancerDNSName"] = stackNames

	return nil
}

func grabSubnetInfo(stub map[string]map[string]interface{}, ec2Client *ec2.Client) error {
	zones := make(map[string]interface{})
	octets := make(map[string]interface{})
	cidrs := make(map[string]interface{})

	stub["SubnetAvailabilityZone"] = zones
	stub["SubnetOctets"] = octets
	stub["SubnetCIDR"] = cidrs

	subnetIds := []string{}
	nameForId := make(map[string]string)
	for name, id := range stub["Subnet"] {
		subnetIds = append(subnetIds, id.(string))
		nameForId[id.(string)] = name
	}

	subnets, err := ec2Client.DescribeSubnets(
		&ec2.DescribeSubnetsParameters{
			Filters: []*ec2.Filter{
				{
					Name:   "subnet-id",
					Values: subnetIds,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	for _, subnet := range subnets.Subnets {
		name := nameForId[subnet.SubnetId]

		_, ipNet, err := net.ParseCIDR(subnet.CidrBlock)
		if err != nil {
			return err
		}

		zones[name] = subnet.AvailabilityZone
		cidrs[name] = subnet.CidrBlock
		octets[name] = []string{
			fmt.Sprintf("%d", ipNet.IP[0]),
			fmt.Sprintf("%d.%d", ipNet.IP[0], ipNet.IP[1]),
			fmt.Sprintf("%d.%d.%d", ipNet.IP[0], ipNet.IP[1], ipNet.IP[2]),
			fmt.Sprintf("%d.%d.%d.%d", ipNet.IP[0], ipNet.IP[1], ipNet.IP[2], ipNet.IP[3]),
		}
	}

	return nil
}
