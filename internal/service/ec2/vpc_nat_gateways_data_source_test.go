package ec2_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccVPCNATGatewaysDataSource_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNATGatewaysDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_nat_gateways.by_vpc_id", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.aws_nat_gateways.by_tags", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.aws_nat_gateways.by_filter", "ids.#", "2"),
					resource.TestCheckResourceAttr("data.aws_nat_gateways.empty", "ids.#", "0"),
				),
			},
		},
	})
}

func testAccNATGatewaysDataSourceConfig(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_vpc" "test1" {
  cidr_block = "172.5.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpc" "test2" {
  cidr_block = "172.6.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_eip" "test1" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_eip" "test2" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test1" {
  vpc_id = aws_vpc.test1.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test2" {
  vpc_id = aws_vpc.test2.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_nat_gateway" "test1" {
  vpc_id = aws_vpc.test1.id

  availability_zone_addresses {
    availability_zone = data.aws_availability_zones.available.names[0]
    allocation_id     = aws_eip.test1.id
  }

  tags = {
    Name     = %[1]q
    OtherTag = "some-value"
  }

  depends_on = [aws_internet_gateway.test1]
}

resource "aws_nat_gateway" "test2" {
  vpc_id = aws_vpc.test2.id

  availability_zone_addresses {
    availability_zone = data.aws_availability_zones.available.names[0]
    allocation_id     = aws_eip.test2.id
  }

  tags = {
    Name     = %[1]q
    OtherTag = "some-other-value"
  }

  depends_on = [aws_internet_gateway.test2]
}

data "aws_nat_gateways" "by_vpc_id" {
  vpc_id = aws_vpc.test2.id

  depends_on = [aws_nat_gateway.test1, aws_nat_gateway.test2]
}

data "aws_nat_gateways" "by_tags" {
  filter {
    name   = "state"
    values = ["available"]
  }

  tags = {
    OtherTag = "some-value"
  }

  depends_on = [aws_nat_gateway.test1, aws_nat_gateway.test2]
}

data "aws_nat_gateways" "by_filter" {
  filter {
    name   = "vpc-id"
    values = [aws_vpc.test1.id, aws_vpc.test2.id]
  }

  depends_on = [aws_nat_gateway.test1, aws_nat_gateway.test2]
}

data "aws_nat_gateways" "empty" {
  vpc_id = aws_vpc.test2.id

  tags = {
    OtherTag = "some-value"
  }

  depends_on = [aws_nat_gateway.test1, aws_nat_gateway.test2]
}
`, rName))
}
