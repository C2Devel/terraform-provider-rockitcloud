package ec2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccVPCNATGatewayDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceNameById := "data.aws_nat_gateway.test_by_id"
	dataSourceNameByVpcId := "data.aws_nat_gateway.test_by_vpc_id"
	dataSourceNameByTags := "data.aws_nat_gateway.test_by_tags"
	resourceName := "aws_nat_gateway.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNATGatewayDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNameById, "availability_mode", resourceName, "availability_mode"),
					resource.TestCheckResourceAttrPair(dataSourceNameById, "connectivity_type", resourceName, "connectivity_type"),
					resource.TestCheckResourceAttrPair(dataSourceNameById, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceNameByVpcId, "vpc_id", resourceName, "vpc_id"),
					resource.TestCheckResourceAttrPair(dataSourceNameByVpcId, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceNameByTags, "tags.Name", resourceName, "tags.Name"),
					resource.TestCheckResourceAttrSet(dataSourceNameById, "state"),
					resource.TestCheckResourceAttr(dataSourceNameById, "nat_gateway_addresses.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(dataSourceNameById, "nat_gateway_addresses.*.allocation_id", "aws_eip.test", "id"),
					resource.TestCheckTypeSetElemAttrPair(dataSourceNameById, "nat_gateway_addresses.*.public_ip", "aws_eip.test", "public_ip"),
					resource.TestCheckResourceAttrSet(dataSourceNameById, "tags.OtherTag"),
				),
			},
		},
	})
}

func testAccNATGatewayDataSourceConfig(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "172.5.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_eip" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_nat_gateway" "test" {
  vpc_id = aws_vpc.test.id

  availability_zone_addresses {
    availability_zone = data.aws_availability_zones.available.names[0]
    allocation_id     = aws_eip.test.id
  }

  tags = {
    Name     = %[1]q
    OtherTag = "some-value"
  }

  depends_on = [aws_internet_gateway.test]
}

data "aws_nat_gateway" "test_by_id" {
  id = aws_nat_gateway.test.id
}

data "aws_nat_gateway" "test_by_vpc_id" {
  vpc_id = aws_nat_gateway.test.vpc_id
}

data "aws_nat_gateway" "test_by_tags" {
  tags = {
    Name = aws_nat_gateway.test.tags["Name"]
  }
}
`, rName))
}
