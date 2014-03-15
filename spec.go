package main

type DeploymentSpec struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	Domain string `yaml:"domain"`

	VPC VPCSpec `yaml:"vpc"`

	InternetGateways []InternetGatewaySpec `yaml:"internet_gateways"`

	Subnets []SubnetSpec `yaml:"subnets"`

	DNS []string `yaml:"dns"`

	SecurityGroups []SecurityGroupSpec `yaml:"security_groups"`

	LoadBalancers []LoadBalancerSpec `yaml:"load_balancers"`

	ElasticIPs []ElasticIPSpec `yaml:"elastic_ips"`
}

type VPCSpec struct {
	CIDR            string `yaml:"cidr"`
	InternetGateway string `yaml:"internet_gateway"`
}

type InternetGatewaySpec struct {
	Name string `yaml:"name"`
}

type SubnetSpec struct {
	Name             string          `yaml:"name"`
	CIDR             string          `yaml:"cidr"`
	AvailabilityZone string          `yaml:"availability_zone"`
	RouteTable       *RouteTableSpec `yaml:"route_table,omitempty"`
	NAT              *SubnetNATSpec  `yaml:"nat,omitempty"`
}

type RouteTableSpec struct {
	InternetGateway *string `yaml:"internet_gateway,omitempty"`
	Instance        *string `yaml:"instance,omitempty"`
}

type SubnetNATSpec struct {
	Name          string `yaml:"name"`
	InstanceType  string `yaml:"type"`
	IP            string `yaml:"ip"`
	KeyPairName   string `yaml:"key_pair_name"`
	SecurityGroup string `yaml:"security_group"`
}

type SecurityGroupSpec struct {
	Name    string                   `yaml:"name"`
	Ingress []SecurityGroupEntrySpec `yaml:"ingress,omitempty"`
	Egress  []SecurityGroupEntrySpec `yaml:"egress,omitempty"`
}

type SecurityGroupEntrySpec struct {
	CIDR     string `yaml:"cidr"`
	Protocol string `yaml:"protocol,omitempty"`
	Ports    string `yaml:"ports,omitempty"`
}

type LoadBalancerSpec struct {
	Name           string                      `yaml:"name"`
	DNSRecord      string                      `yaml:"dns_record"`
	Listeners      []LoadBalancerListenerSpec  `yaml:"listeners"`
	HealthCheck    LoadBalancerHealthCheckSpec `yaml:"health_check"`
	Subnets        []string                    `yaml:"subnets"`
	SecurityGroups []string                    `yaml:"security_groups"`
}

type LoadBalancerListenerSpec struct {
	Protocol            string  `yaml:"protocol"`
	Port                uint16  `yaml:"port"`
	DestinationProtocol *string `yaml:"destination_protocol,omitempty"`
	DestinationPort     *uint16 `yaml:"destination_port,omitempty"`
}

type LoadBalancerHealthCheckSpec struct {
	Target             LoadBalancerHealthCheckTargetSpec `yaml:"target"`
	Timeout            int                               `yaml:"timeout"`
	Interval           int                               `yaml:"interval"`
	HealthyThreshold   int                               `yaml:"healthy_threshold"`
	UnhealthyThreshold int                               `yaml:"unhealthy_threshold"`
}

type LoadBalancerHealthCheckTargetSpec struct {
	Type string `yaml:"type"`
	Port uint16 `yaml:"port"`
}

type ElasticIPSpec struct {
	Name      string `yaml:"name"`
	DNSRecord string `yaml:"dns_record"`
}
