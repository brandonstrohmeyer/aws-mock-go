package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Mocked EC2 client and methods. Each method is setup to return a pointer to the relevant struct attribute.
// This allows different return values to be set for each individual instantiation of the mocked client.
type mockedEC2 struct {
	ec2iface.EC2API
	DescribeVpcsOutput ec2.DescribeVpcsOutput
}

func (m mockedEC2) DescribeVpcs(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return &m.DescribeVpcsOutput, nil
}

// Define each output to mock as a return value.
var (
	// DescribeVpcsOutput - Single CIDR Block
	describeVpcsOutputSingleAssociation = ec2.DescribeVpcsOutput{
		Vpcs: []*ec2.Vpc{
			{
				CidrBlock: aws.String("172.29.240.0/20"),
				VpcId:     aws.String("vpc-0e3533b1b004a32f1"),
				CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{
					{
						AssociationId: aws.String("vpc-cidr-assoc-0fac8e805efcc0f21"),
						CidrBlock:     aws.String("172.29.240.0/20"),
					},
				},
			},
		},
	}

	// DescribeVpcsOutput - Multiple CIDR Blocks
	describeVpcsOutputMultipleAssociation = ec2.DescribeVpcsOutput{
		Vpcs: []*ec2.Vpc{
			{
				CidrBlock: aws.String("172.30.240.0/20"),
				VpcId:     aws.String("vpc-1f3533b1b004b32f0"),
				CidrBlockAssociationSet: []*ec2.VpcCidrBlockAssociation{
					{
						AssociationId: aws.String("vpc-cidr-assoc-0fac8e105ffcc0f21"),
						CidrBlock:     aws.String("172.30.240.0/20"),
					},
					{
						AssociationId: aws.String("vpc-cidr-assoc-1gac8e305efcc0f22"),
						CidrBlock:     aws.String("172.31.240.0/20"),
					},
				},
			},
		},
	}
)

// Unit test for ListVpcCidrBlocks func
func TestListVpcCidrBlocks(t *testing.T) {
	// Define each test case as a struct.
	// There are two specific cases tested for here:
	//
	// 1. A VPC with a single CIDR block attached
	// 2. A VPC with multiple CIDR blocks attached
	cases := []struct {
		Name     string
		Resp     ec2.DescribeVpcsOutput
		Expected []string
	}{
		{
			Name:     "SingleCidrAssociation",
			Resp:     describeVpcsOutputSingleAssociation,
			Expected: []string{"172.29.240.0/20"},
		},
		{
			Name:     "MultipleCidrAssociation",
			Resp:     describeVpcsOutputMultipleAssociation,
			Expected: []string{"172.30.240.0/20", "172.31.240.0/20"},
		},
	}

	// Iterate through the list of test cases defined above
	for _, c := range cases {

		// Create a sub-test for each case
		t.Run(c.Name, func(t *testing.T) {

			// Run the test using the mocked EC2 client and specified return value
			cidrs, err := ListVpcCidrBlocks(
				mockedEC2{DescribeVpcsOutput: c.Resp},
				"vpc-0e3533b1b004a32f1")

			// Test that the returned string matches what we expect
			assert.Equal(t, c.Expected, cidrs)

			// Test that no error was returned
			assert.NoError(t, err)

		})
	}
}
