package models

import (
	"encoding/json"
	"strings"
)

type Hash map[string]interface{}

func Ref(i interface{}) Hash {
	return Hash{"Ref": i}
}

type Template struct {
	AWSTemplateFormatVersion string
	Description              string
	Resources                Resources
	Mappings                 Hash
}

type Tag struct {
	Key   string `json:"Key,omitempty"`
	Value string `json:"Value,omitempty"`
}

type Resources map[string]Resource

type Resource interface {
	Type() string
	DependsOn() string
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

		if v.DependsOn() != "" {
			p["DependsOn"] = v.DependsOn()
		}

		vj, e := json.Marshal(v)
		if e != nil {
			return nil, e
		}

		if string(vj) != "{}" {
			p["Properties"] = v
		}

		pj, e := json.Marshal(p)
		if e != nil {
			return nil, e
		}

		lines = append(lines, string(kj)+": "+string(pj))
	}
	return []byte("{" + strings.Join(lines, ",\n") + "}"), nil
}

type InternetGateway struct {
}

func (InternetGateway) Type() string {
	return "AWS::EC2::InternetGateway"
}

func (InternetGateway) DependsOn() string {
	return ""
}

type VPC struct {
	CidrBlock        interface{} `json:"CidrBlock,omitempty"`
	EnableDnsSupport interface{} `json:"EnableDnsSupport,omitempty"`
	InstanceTenancy  interface{} `json:"InstanceTenancy,omitempty"`
	Tags             interface{} `json:"Tags,omitempty"`
}

func (VPC) Type() string {
	return "AWS::EC2::VPC"
}

func (VPC) DependsOn() string {
	return ""
}

type VPCGatewayAttachment struct {
	InternetGatewayId interface{} `json:"InternetGatewayId,omitempty"`
	VpcId             interface{} `json:"VpcId,omitempty"`
}

func (VPCGatewayAttachment) Type() string {
	return "AWS::EC2::VPCGatewayAttachment"
}

func (VPCGatewayAttachment) DependsOn() string {
	return ""
}

type RouteTable struct {
	VpcId interface{} `json:"VpcId,omitempty"`
}

func (RouteTable) Type() string {
	return "AWS::EC2::RouteTable"
}

func (RouteTable) DependsOn() string {
	return ""
}

type SecurityGroup struct {
	GroupDescription     interface{} `json:"GroupDescription,omitempty"`
	SecurityGroupEgress  interface{} `json:"SecurityGroupEgress,omitempty"`
	SecurityGroupIngress interface{} `json:"SecurityGroupIngress,omitempty"`
	VpcId                interface{} `json:"VpcId,omitempty"`
	Tags                 interface{} `json:"Tags,omitempty"`
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

func (SecurityGroup) Type() string {
	return "AWS::EC2::SecurityGroup"
}

func (SecurityGroup) DependsOn() string {
	return ""
}

type LoadBalancer struct {
	Subnets        interface{} `json:"Subnets,omitempty"`
	HealthCheck    interface{} `json:"HealthCheck,omitempty"`
	Instances      interface{} `json:"Instances,omitempty"`
	SecurityGroups interface{} `json:"SecurityGroups,omitempty"`
	Listeners      interface{} `json:"Listeners,omitempty"`
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
	SSLCertificateId string `json:"SSLCertificateId,omitempty"`
}

func (LoadBalancer) Type() string {
	return "AWS::ElasticLoadBalancing::LoadBalancer"
}

func (LoadBalancer) DependsOn() string {
	return ""
}

type Route struct {
	DestinationCidrBlock interface{} `json:"DestinationCidrBlock,omitempty"`
	GatewayId            interface{} `json:"GatewayId,omitempty"`
	RouteTableId         interface{} `json:"RouteTableId,omitempty"`
	InstanceId           interface{} `json:"InstanceId,omitempty"`
	NetworkInterfaceId   interface{} `json:"NetworkInterfaceId,omitempty"`
	Depends              string      `json:"-"`
}

func (Route) Type() string {
	return "AWS::EC2::Route"
}

func (route Route) DependsOn() string {
	return route.Depends
}

type Subnet struct {
	AvailabilityZone interface{} `json:"AvailabilityZone,omitempty"`
	CidrBlock        interface{} `json:"CidrBlock,omitempty"`
	VpcId            interface{} `json:"VpcId,omitempty"`
	Tags             interface{} `json:"Tags,omitempty"`
}

func (Subnet) Type() string {
	return "AWS::EC2::Subnet"
}

func (Subnet) DependsOn() string {
	return ""
}

type SubnetRouteTableAssociation struct {
	RouteTableId interface{} `json:"RouteTableId,omitempty"`
	SubnetId     interface{} `json:"SubnetId,omitempty"`
}

func (SubnetRouteTableAssociation) Type() string {
	return "AWS::EC2::SubnetRouteTableAssociation"
}

func (SubnetRouteTableAssociation) DependsOn() string {
	return ""
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

func (Instance) Type() string {
	return "AWS::EC2::Instance"
}

func (Instance) DependsOn() string {
	return ""
}

type DHCPOptions struct {
	DomainName         interface{} `json:"DomainName,omitempty"`
	DomainNameServers  interface{} `json:"DomainNameServers,omitempty"`
	NetbiosNameServers interface{} `json:"NetbiosNameServers,omitempty"`
	NetbiosNodeType    interface{} `json:"NetbiosNodeType,omitempty"`
	NtpServers         interface{} `json:"NtpServers,omitempty"`
	Tags               interface{} `json:"Tags,omitempty"`
}

func (DHCPOptions) Type() string {
	return "AWS::EC2::DHCPOptions"
}

func (DHCPOptions) DependsOn() string {
	return ""
}

type VPCDHCPOptionsAssociation struct {
	VpcId         interface{} `json:"VpcId,omitempty"`
	DhcpOptionsId interface{} `json:"DhcpOptionsId,omitempty"`
}

func (VPCDHCPOptionsAssociation) Type() string {
	return "AWS::EC2::VPCDHCPOptionsAssociation"
}

func (VPCDHCPOptionsAssociation) DependsOn() string {
	return ""
}

type EIP struct {
	Domain     interface{} `json:"Domain,omitempty"`
	InstanceId interface{} `json:"InstanceId,omitempty"`
	Depends    string      `json:"-"`
}

func (EIP) Type() string {
	return "AWS::EC2::EIP"
}

func (eip EIP) DependsOn() string {
	return eip.Depends
}

type RecordSetGroup struct {
	HostedZoneName interface{} `json:"HostedZoneName,omitempty"`
	RecordSets     interface{} `json:"RecordSets,omitempty"`
}

func (RecordSetGroup) Type() string {
	return "AWS::Route53::RecordSetGroup"
}

func (RecordSetGroup) DependsOn() string {
	return ""
}

type RecordSet struct {
	Name        interface{} `json:"Name,omitempty"`
	Type        interface{} `json:"Type,omitempty"`
	AliasTarget interface{} `json:"AliasTarget,omitempty"`
}

type RecordSetAliasTarget struct {
	HostedZoneId interface{} `json:"HostedZoneId,omitempty"`
	DNSName      interface{} `json:"DNSName,omitempty"`
}

type Bucket struct {
	BucketName string `json:"BucketName,omitempty"`
}

func (Bucket) Type() string {
	return "AWS::S3::Bucket"
}

func (Bucket) DependsOn() string {
	return ""
}
