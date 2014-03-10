package main

import (
	"encoding/json"
	"strings"
)

// base types
type Hash map[string]interface{}
type List []interface{}

// special types

type Template struct {
	AWSTemplateFormatVersion string
	Description              string
	Resources                Resources
	Mappings                 Hash
}

type Typer interface {
	Type() string
}

type InternetGateway struct {
}

func (g InternetGateway) Type() string {
	return "AWS::EC2::InternetGateway"
}

type Vpc struct {
	CidrBlock        interface{} `json:"CidrBlock,omitempty"`
	EnableDnsSupport interface{} `json:"EnableDnsSupport,omitempty"`
	InstanceTenancy  interface{} `json:"InstanceTenancy,omitempty"`
	Tags             interface{} `json:"Tags,omitempty"`
}

type VPCGatewayAttachment struct {
	InternetGatewayId interface{} `json:"InternetGatewayId,omitempty"`
	VpcId             interface{} `json:"VpcId,omitempty"`
}

type RouteTable struct {
	VpcId interface{} `json:"VpcId,omitempty"`
}

type SecurityGroup struct {
	GroupDescription     interface{} `json:"GroupDescription,omitempty"`
	SecurityGroupEgress  interface{} `json:"SecurityGroupEgress,omitempty"`
	SecurityGroupIngress interface{} `json:"SecurityGroupIngress,omitempty"`
	VpcId                interface{} `json:"VpcId,omitempty"`
	Tags                 interface{} `json:"Tags,omitempty"`
}

func (s SecurityGroup) Type() string {
	return "AWS::EC2::SecurityGroup"
}

type SecurityGroupEgress struct {
	CidrIp                     interface{} `json:"CidrIp,omitempty"`
	FromPort                   interface{} `json:"FromPort,omitempty"`
	IpProtocol                 interface{} `json:"IpProtocol,omitempty"`
	DestinationSecurityGroupId interface{} `json:"DestinationSecurityGroupId,omitempty"`
	ToPort                     interface{} `json:"ToPort,omitempty"`
}

type SecurityGroupIngress struct {
	CidrIp                     interface{} `json:"CidrIp,omitempty"`
	FromPort                   interface{} `json:"FromPort,omitempty"`
	IpProtocol                 interface{} `json:"IpProtocol,omitempty"`
	SourceSecurityGroupId      interface{} `json:"SourceSecurityGroupId,omitempty"`
	SourceSecurityGroupName    interface{} `json:"SourceSecurityGroupName,omitempty"`
	SourceSecurityGroupOwnerId interface{} `json:"SourceSecurityGroupOwnerId,omitempty"`
	ToPort                     interface{} `json:"ToPort,omitempty"`
}

func (r RouteTable) Type() string {
	return "AWS::EC2::RouteTable"
}

func (e VPCGatewayAttachment) Type() string {
	return "AWS::EC2::VPCGatewayAttachment"
}

type LoadBalancer struct {
	Subnets        interface{} `json:"Subnets,omitempty"`
	HealthCheck    interface{} `json:"HealthCheck,omitempty"`
	Instances      interface{} `json:"Instances,omitempty"`
	SecurityGroups interface{} `json:"SecurityGroups,omitempty"`
	Listeners      interface{} `json:"Listeners,omitempty"`
}

func (r LoadBalancer) Type() string {
	return "AWS::ElasticLoadBalancing::LoadBalancer"
}

type LoadBalancerHealthCheck struct {
	HealthyThreshold   string `json:"HealthyThreshold,omitempty"`
	Interval           string `json:"Interval,omitempty"`
	Target             string `json:"Target,omitempty"`
	Timeout            string `json:"Timeout,omitempty"`
	UnhealthyThreshold string `json:"UnhealthyThreshold,omitempty"`
}

type LoadBalancerListener struct {
	InstancePort     string `json:"InstancePort,omitempty"`
	LoadBalancerPort string `json:"LoadBalancerPort,omitempty"`
	Protocol         string `json:"Protocol,omitempty"`
	InstanceProtocol string `json:"InstanceProtocol,omitempty"`
}

type Route struct {
	DestinationCidrBlock interface{} `json:"DestinationCidrBlock,omitempty"`
	GatewayId            interface{} `json:"GatewayId,omitempty"`
	RouteTableId         interface{} `json:"RouteTableId,omitempty"`
	InstanceId           interface{} `json:"InstanceId,omitempty"`
	NetworkInterfaceId   interface{} `json:"NetworkInterfaceId,omitempty"`
}

func (r Route) Type() string {
	return "AWS::EC2::Route"
}

type Subnet struct {
	AvailabilityZone interface{} `json:"AvailabilityZone,omitempty"`
	CidrBlock        interface{} `json:"CidrBlock,omitempty"`
	VpcId            interface{} `json:"VpcId,omitempty"`
	Tags             interface{} `json:"Tags,omitempty"`
}

type SubnetRouteTableAssociation struct {
	RouteTableId interface{} `json:"RouteTableId,omitempty"`
	SubnetId     interface{} `json:"SubnetId,omitempty"`
}

func (s SubnetRouteTableAssociation) Type() string {
	return "AWS::EC2::SubnetRouteTableAssociation"
}

func (s Subnet) Type() string {
	return "AWS::EC2::Subnet"
}

func (r Resources) MarshalJSON() ([]byte, error) {
	lines := []string{}
	for k, v := range r {
		kj, e := json.Marshal(k)
		if e != nil {
			return nil, e
		}
		p := map[string]interface{}{
			"Type": v.Type(),
		}
		vj, e := json.Marshal(v)
		if e != nil {
			return nil, e
		}
		if string(vj) != "{}" {
			p["Properties"] = v
		}
		pj, e := json.Marshal(p)
		lines = append(lines, string(kj)+": "+string(pj))
	}
	return []byte("{" + strings.Join(lines, ",\n") + "}"), nil
}

func (v Vpc) Type() string {
	return "AWS::EC2::VPC"
}

type Instance struct {
	AvailabilityZone      interface{} `json:"AvailabilityZone,omitempty"`
	BlockDeviceMappings   interface{} `json:"BlockDeviceMappings,omitempty"`
	DisableApiTermination interface{} `json:"DisableApiTermination,omitempty"`
	EbsOptimized          interface{} `json:"EbsOptimized,omitempty"`
	IamInstanceProfile    interface{} `json:"IamInstanceProfile,omitempty"`
	ImageId               interface{} `json:"ImageId,omitempty"`
	InstanceType          interface{} `json:"InstanceType,omitempty"`
	KernelId              interface{} `json:"KernelId,omitempty"`
	KeyName               interface{} `json:"KeyName,omitempty"`
	Monitoring            interface{} `json:"Monitoring,omitempty"`
	NetworkInterfaces     interface{} `json:"NetworkInterfaces,omitempty"`
	PlacementGroupName    interface{} `json:"PlacementGroupName,omitempty"`
	PrivateIpAddress      interface{} `json:"PrivateIpAddress,omitempty"`
	RamdiskId             interface{} `json:"RamdiskId,omitempty"`
	SecurityGroupIds      interface{} `json:"SecurityGroupIds,omitempty"`
	SecurityGroups        interface{} `json:"SecurityGroups,omitempty"`
	SourceDestCheck       interface{} `json:"SourceDestCheck,omitempty"`
	SubnetId              interface{} `json:"SubnetId,omitempty"`
	Tags                  interface{} `json:"Tags,omitempty"`
	Tenancy               interface{} `json:"Tenancy,omitempty"`
	UserData              interface{} `json:"UserData,omitempty"`
	Volumes               interface{} `json:"Volumes,omitempty"`
}

func (i Instance) Type() string {
	return "AWS::EC2::Instance"
}

type NetworkInterface struct {
	AssociatePublicIpAddress       interface{} `json:"AssociatePublicIpAddress,omitempty"`
	DeleteOnTermination            interface{} `json:"DeleteOnTermination,omitempty"`
	Description                    interface{} `json:"Description,omitempty"`
	DeviceIndex                    interface{} `json:"DeviceIndex,omitempty"`
	GroupSet                       interface{} `json:"GroupSet,omitempty"`
	NetworkInterfaceId             interface{} `json:"NetworkInterfaceId,omitempty"`
	PrivateIpAddress               interface{} `json:"PrivateIpAddress,omitempty"`
	PrivateIpAddresses             interface{} `json:"PrivateIpAddresses,omitempty"`
	SecondaryPrivateIpAddressCount interface{} `json:"SecondaryPrivateIpAddressCount,omitempty"`
	SubnetId                       interface{} `json:"SubnetId,omitempty"`
}

type DHCPOptions struct {
	DomainName         interface{} `json:"DomainName,omitempty"`
	DomainNameServers  interface{} `json:"DomainNameServers,omitempty"`
	NetbiosNameServers interface{} `json:"NetbiosNameServers,omitempty"`
	NetbiosNodeType    interface{} `json:"NetbiosNodeType,omitempty"`
	NtpServers         interface{} `json:"NtpServers,omitempty"`
	Tags               interface{} `json:"Tags,omitempty"`
}

func (d DHCPOptions) Type() string {
	return "AWS::EC2::DHCPOptions"
}

type VPCDHCPOptionsAssociation struct {
	VpcId         interface{} `json:"VpcId,omitempty"`
	DhcpOptionsId interface{} `json:"DhcpOptionsId,omitempty"`
}

func (v VPCDHCPOptionsAssociation) Type() string {
	return "AWS::EC2::VPCDHCPOptionsAssociation"
}

type Tag struct {
	Key   string `json:"Key,omitempty"`
	Value string `json:"Value,omitempty"`
}

type Resources map[string]Typer

func ref(i interface{}) Hash {
	return Hash{"Ref": i}
}
