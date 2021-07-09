package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"log"
)

// List all CIDR blocks attached to a VPC
func ListVpcCidrBlocks(svc ec2iface.EC2API, vpcId string) ([]string, error) {
	var c []string
	result, err := svc.DescribeVpcs(
		&ec2.DescribeVpcsInput{
			VpcIds: aws.StringSlice([]string{vpcId}),
		},
	)
	if err != nil {
		return c, fmt.Errorf("Unable to describe VPCs: %v", err)
	}

	// Iterate through the CIDR blocks associated with the returned VPC.
	// Since we're passing a "VpcIds" filter, only a single VPC will be returned in the slice.
	for _, s := range result.Vpcs[0].CidrBlockAssociationSet {
		c = append(c, *s.CidrBlock)
	}

	return c, nil
}

func main() {
	// Create AWS API Session
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String("us-east-1")},
	})
	if err != nil {
		log.Fatal("Error creating AWS session:", err.Error())
	}

	// Instantiate new EC2 client
	svc := ec2.New(sess)

	// Get a list of CIDR blocks attached to the VPC
	CidrBlocks, err := ListVpcCidrBlocks(svc, "vpc-09ce5efd93f89f655")
	if err != nil {
		log.Fatal("Error listing VPC cidr blocks:", err.Error())
	}
	fmt.Println(CidrBlocks)
}
