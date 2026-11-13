package ec2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccVPNSiteGatewayDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceNameByID := "data.aws_vpn_gateway.test_by_id"
	dataSourceNameByVPCID := "data.aws_vpn_gateway.test_by_vpc_id"
	dataSourceNameByAttachedVPCID := "data.aws_vpn_gateway.test_by_attached_vpc_id"
	dataSourceNameByFilter := "data.aws_vpn_gateway.test_by_filter"
	vpcResourceName := "aws_vpc.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNGatewayDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceNameByID, "id", regexp.MustCompile(`^vgw-.+`)),
					acctest.MatchResourceAttrRegionalARN(dataSourceNameByID, "arn", "ec2", regexp.MustCompile(`vpn-gateway/vgw-.+`)),
					resource.TestCheckResourceAttrPair(dataSourceNameByID, "attached_vpc_id", vpcResourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceNameByID, "state", "available"),
					resource.TestCheckResourceAttrPair(dataSourceNameByVPCID, "id", dataSourceNameByID, "id"),
					resource.TestCheckResourceAttrPair(dataSourceNameByAttachedVPCID, "id", dataSourceNameByID, "id"),
					resource.TestCheckResourceAttrPair(dataSourceNameByFilter, "id", dataSourceNameByID, "id"),
				),
			},
		},
	})
}

// VPN gateway is created together with VPC.
func testAccVPNGatewayDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

data "aws_vpn_gateway" "test_by_id" {
  id = replace(aws_vpc.test.id, "vpc-", "vgw-")
}

data "aws_vpn_gateway" "test_by_vpc_id" {
  id = aws_vpc.test.id
}

data "aws_vpn_gateway" "test_by_attached_vpc_id" {
  attached_vpc_id = aws_vpc.test.id
  state           = "available"
}

data "aws_vpn_gateway" "test_by_filter" {
  filter {
    name   = "attachment.vpc-id"
    values = [aws_vpc.test.id]
  }
}
`, rName)
}
