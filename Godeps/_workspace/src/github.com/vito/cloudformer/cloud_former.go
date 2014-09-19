package cloudformer

import (
	"net"
	"time"
)

type CloudFormer interface {
	InternetGateway(name string) InternetGateway
	VPC(name string) VPC
	ElasticIP(name string) ElasticIP
	LoadBalancer(name string) LoadBalancer
	Bucket(name string) Bucket
}

type InternetGateway interface{}

type DHCPOptions struct {
	DomainNameServers []string
}

type VPC interface {
	Network(*net.IPNet)

	AttachInternetGateway(InternetGateway)
	AssociateDHCPOptions(DHCPOptions)

	Subnet(name string) Subnet
	SecurityGroup(name string) SecurityGroup
}

type Subnet interface {
	Network(*net.IPNet)
	AvailabilityZone(string)

	Instance(name string) Instance
	RouteTable() RouteTable
}

type SecurityGroup interface {
	Ingress(ProtocolType, *net.IPNet, uint16, uint16)
	Egress(ProtocolType, *net.IPNet, uint16, uint16)
}

type RouteTable interface {
	InternetGateway(InternetGateway)
	Instance(Instance)
}

type Instance interface {
	Type(string)
	Image(string)
	PrivateIP(net.IP)
	KeyPair(string)
	SecurityGroup(SecurityGroup)
	SourceDestCheck(bool)
}

type ElasticIP interface {
	Domain(name string)

	AttachTo(Instance)
}

type LoadBalancer interface {
	Listener(ProtocolType, uint16, ProtocolType, uint16, string)
	HealthCheck(HealthCheck)
	Subnet(Subnet)
	SecurityGroup(SecurityGroup)
	RecordSet(name, zone string)
	Scheme(scheme string)
}

type HealthCheck struct {
	Protocol           ProtocolType
	Port               uint16
	Path               string
	Timeout            time.Duration
	Interval           time.Duration
	HealthyThreshold   int
	UnhealthyThreshold int
}

type ProtocolType string

const TCP = ProtocolType("tcp")
const UDP = ProtocolType("udp")

type Bucket interface {
	Name(string)
}
