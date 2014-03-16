package main_test

import (
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vito/cmdtest"
	. "github.com/vito/cmdtest/matchers"
)

var _ = Describe("Boosh", func() {
	var booshPath string

	BeforeEach(func() {
		var err error

		booshPath, err = cmdtest.Build("github.com/vito/boosh")
		Ω(err).ShouldNot(HaveOccurred())
	})

	generating := func() *cmdtest.Session {
		session, err := cmdtest.Start(exec.Command(booshPath, "generate"))
		Ω(err).ShouldNot(HaveOccurred())

		return session
	}

	It("generates a correct deployment", func() {
		session := generating()

		session.Stdin.Write([]byte(`---
description: MicroBOSH

elastic_ips:
  - name: Micro

vpc:
  cidr: 10.11.0.0/16
  internet_gateway: VPCGateway

internet_gateways:
  - name: VPCGateway

subnets:
  - name: BOSHZ1
    cidr: 10.11.0.0/24
    availability_zone: us-east-1a
    route_table:
      internet_gateway: VPCGateway

dns:
  - 10.11.0.2
  - 10.11.0.6

security_groups:
  - name: BOSH
    ingress:
      - protocol: tcp
        ports: 0-65535
        cidr: 0.0.0.0/0
      - protocol: udp
        ports: 0-65535
        cidr: 0.0.0.0/0
`))

		session.Stdin.Close()

		Ω(session).Should(ExitWithTimeout(0, 1*time.Second))

		Ω(session.FullOutput()).Should(MatchJSON(`{
  "Mappings": {},
  "Resources": {
    "BOSHZ1Route": {
      "Type": "AWS::EC2::Route",
      "Properties": {
        "RouteTableId": {
          "Ref": "BOSHZ1RouteTable"
        },
        "GatewayId": {
          "Ref": "VPCGatewayInternetGateway"
        },
        "DestinationCidrBlock": "0.0.0.0/0"
      },
      "DependsOn": "VPCGatewayInternetGateway"
    },
    "BOSHZ1SubnetRouteTableAssociation": {
      "Type": "AWS::EC2::SubnetRouteTableAssociation",
      "Properties": {
        "SubnetId": {
          "Ref": "BOSHZ1Subnet"
        },
        "RouteTableId": {
          "Ref": "BOSHZ1RouteTable"
        }
      }
    },
    "BOSHZ1RouteTable": {
      "Type": "AWS::EC2::RouteTable",
      "Properties": {
        "VpcId": {
          "Ref": "VPC"
        }
      }
    },
    "VPC": {
      "Type": "AWS::EC2::VPC",
      "Properties": {
        "CidrBlock": "10.11.0.0/16"
      }
    },
    "DHCPOptions": {
      "Type": "AWS::EC2::DHCPOptions",
      "Properties": {
        "DomainNameServers": [
          "10.11.0.2",
          "10.11.0.6"
        ]
      }
    },
    "MicroEIP": {
      "Type": "AWS::EC2::EIP",
      "Properties": {
        "Domain": "vpc"
      }
    },
    "VPCDHCPOptionsAssociation": {
      "Type": "AWS::EC2::VPCDHCPOptionsAssociation",
      "Properties": {
        "DhcpOptionsId": {
          "Ref": "DHCPOptions"
        },
        "VpcId": {
          "Ref": "VPC"
        }
      }
    },
    "VPCGatewayInternetGateway": {
      "Type": "AWS::EC2::InternetGateway"
    },
    "VPCGatewayAttachment": {
      "Type": "AWS::EC2::VPCGatewayAttachment",
      "Properties": {
        "VpcId": {
          "Ref": "VPC"
        },
        "InternetGatewayId": {
          "Ref": "VPCGatewayInternetGateway"
        }
      }
    },
    "BOSHSecurityGroup": {
      "Type": "AWS::EC2::SecurityGroup",
      "Properties": {
        "VpcId": {
          "Ref": "VPC"
        },
        "SecurityGroupIngress": [
          {
            "ToPort": "65535",
            "IpProtocol": "tcp",
            "FromPort": "0",
            "CidrIp": "0.0.0.0/0"
          },
          {
            "ToPort": "65535",
            "IpProtocol": "udp",
            "FromPort": "0",
            "CidrIp": "0.0.0.0/0"
          }
        ],
        "SecurityGroupEgress": [],
        "GroupDescription": "BOSH"
      }
    },
    "BOSHZ1Subnet": {
      "Type": "AWS::EC2::Subnet",
      "Properties": {
        "Tags": [
          {
            "Value": "BOSHZ1",
            "Key": "Name"
          }
        ],
        "VpcId": {
          "Ref": "VPC"
        },
        "CidrBlock": "10.11.0.0/24",
        "AvailabilityZone": "us-east-1a"
      }
    }
  },
  "Description": "MicroBOSH",
  "AWSTemplateFormatVersion": "2010-09-09"
}`))
	})
})
