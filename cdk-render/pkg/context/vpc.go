package context

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type VpcNetworkPlugin struct {
	vpcId string
}

type vpcNetwork struct {
	VpcId           string        `json:"vpcId"`
	VpcCidrBlock    string        `json:"vpcCidrBlock"`
	AvailbiltyZones []string      `json:"availabilityZones"`
	SubnetGroups    []subnetGroup `json:"subnetGroups"`
}

type subnetGroup struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Subnets []subnet `json:"subnets"`
}

type subnet struct {
	SubnetId         string `json:"subnetId"`
	Cidr             string `json:"cidr"`
	AvailabilityZone string `json:"availabilityZone"`
	RouteTableId     string `json:"routeTableId"`
}

func NewVpcPlugin(vpcId string) *VpcNetworkPlugin {
	return &VpcNetworkPlugin{
		vpcId: vpcId,
	}
}

func (v *VpcNetworkPlugin) Render(ctx *CdkContext) (map[string]any, error) {
	key := fmt.Sprintf("vpc-provider:account=%s:filter.vpc-id=%s:region=%s:returnAsymmetricSubnets=true", ctx.AwsMeta.Account, v.vpcId, ctx.AwsMeta.Region)

	d, err := getVpcInfo(v.vpcId, ctx.Context, ctx.Ec2Client)
	return map[string]any{key: d}, err
}

func getVpcInfo(vpcId string, ctx context.Context, c *ec2.Client) (*vpcNetwork, error) {
	resp := &vpcNetwork{
		VpcId:        vpcId,
		SubnetGroups: []subnetGroup{},
	}

	vpc, err := c.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
		VpcIds: []string{vpcId},
	})
	if err != nil {
		return nil, err
	}

	if len(vpc.Vpcs) > 0 {
		resp.VpcCidrBlock = *vpc.Vpcs[0].CidrBlock
		resp.AvailbiltyZones = []string{}
	}

	subnets, err := c.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcId},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	routeTables, err := c.DescribeRouteTables(ctx, &ec2.DescribeRouteTablesInput{})

	resp.SubnetGroups = processSubnets(subnets, routeTables)
	return resp, err
}

func processSubnets(subnets *ec2.DescribeSubnetsOutput, routeTables *ec2.DescribeRouteTablesOutput) []subnetGroup {
	subMap := map[string]subnetGroup{}
	for _, s := range subnets.Subnets {
		// CDK created subnets are tagged...
		// TODO: Follow similar logic as CDK for imported
		name, ok := getTag("aws-cdk:subnet-name", s.Tags)
		if !ok {
			name = ""
		}

		subnetType, ok := getTag("aws-cdk:subnet-type", s.Tags)
		if !ok {
			subnetType = ""
			if *s.MapPublicIpOnLaunch {
				subnetType = "Public"
			}
			subnetType = ""
		}

		// CDK behavior
		if name == "" {
			name = subnetType
		}

		if _, ok := subMap[name]; !ok {
			subMap[name] = subnetGroup{
				Name:    name,
				Type:    subnetType,
				Subnets: []subnet{},
			}
		}

		grp := subMap[name]
		grp.Subnets = append(subMap[name].Subnets, subnet{
			SubnetId:         *s.SubnetId,
			Cidr:             *s.CidrBlock,
			AvailabilityZone: *s.AvailabilityZone,
			RouteTableId:     routeTableIdForSubnet(*s.SubnetId, routeTables.RouteTables),
		})
		subMap[name] = grp
	}

	resp := []subnetGroup{}
	for _, v := range subMap {
		resp = append(resp, v)
	}
	return resp
}

func getTag(key string, tags []types.Tag) (string, bool) {
	for _, t := range tags {
		if *t.Key == key {
			return *t.Value, true
		}
	}
	return "", false
}

func routeTableIdForSubnet(subnetId string, rtbls []types.RouteTable) string {
	for _, rt := range rtbls {
		for _, asc := range rt.Associations {
			if asc.SubnetId != nil && subnetId == *asc.SubnetId {
				return *rt.RouteTableId
			}
		}
	}
	return ""
}
