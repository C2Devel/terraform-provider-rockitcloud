package ec2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccVPCNATGateway_basic(t *testing.T) {
	var natGateway ec2.NatGateway
	var route ec2.Route
	resourceName := "aws_nat_gateway.test"
	routeResourceName := "aws_route.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	destinationCidr := "10.3.0.0/16"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckNATGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNATGatewayConfigRoute(rName, destinationCidr),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNATGatewayExists(resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "availability_mode", "regional"),
					resource.TestCheckResourceAttr(resourceName, "auto_provision_zones", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "connectivity_type", "public"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "aws_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "availability_zone_addresses.#", "1"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "availability_zone_addresses.*.allocation_id", "aws_eip.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_addresses.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "nat_gateway_addresses.*", map[string]string{
						"is_primary": "true",
					}),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "nat_gateway_addresses.*.allocation_id", "aws_eip.test", "id"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "nat_gateway_addresses.*.public_ip", "aws_eip.test", "public_ip"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					// route targets NAT gateway
					testAccCheckRouteExists(routeResourceName, &route),
					resource.TestCheckResourceAttrPair(routeResourceName, "nat_gateway_id", resourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVPCNATGateway_disappears(t *testing.T) {
	var natGateway ec2.NatGateway
	resourceName := "aws_nat_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckNATGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNATGatewayConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNATGatewayExists(resourceName, &natGateway),
					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceNATGateway(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVPCNATGateway_autoAllocate(t *testing.T) {
	var natGateway ec2.NatGateway
	resourceName := "aws_nat_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckNATGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNATGatewayConfigAutoAllocate(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNATGatewayExists(resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "availability_mode", "regional"),
					resource.TestCheckResourceAttr(resourceName, "auto_provision_zones", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "connectivity_type", "public"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", "aws_vpc.test", "id"),
					// automatically allocated addresses populated only for VPC with running AZ
					resource.TestCheckResourceAttr(resourceName, "availability_zone_addresses.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_addresses.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVPCNATGateway_availabilityZoneAddresses(t *testing.T) {
	var natGateway1, natGateway2, natGateway3 ec2.NatGateway
	resourceName := "aws_nat_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t); testAccPreCheckAvailabilityZonesAtLeast(t, 2) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckNATGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNATGatewayConfigAvailabilityZoneAddresses1(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNATGatewayExists(resourceName, &natGateway1),
					resource.TestCheckResourceAttr(resourceName, "availability_zone_addresses.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_addresses.#", "1"),
				),
			},
			{
				// Associate EIP with second availability zone
				Config: testAccNATGatewayConfigAvailabilityZoneAddresses2(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNATGatewayExists(resourceName, &natGateway2),
					testAccCheckNATGatewayNotRecreated(&natGateway1, &natGateway2),
					resource.TestCheckResourceAttr(resourceName, "availability_zone_addresses.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "nat_gateway_addresses.#", "2"),
				),
			},
			{
				// Swap EIP in first availability zone
				Config: testAccNATGatewayConfigAvailabilityZoneAddresses3(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNATGatewayExists(resourceName, &natGateway3),
					testAccCheckNATGatewayNotRecreated(&natGateway2, &natGateway3),
					resource.TestCheckResourceAttr(resourceName, "availability_zone_addresses.#", "2"),
					resource.TestCheckTypeSetElemAttrPair(resourceName, "availability_zone_addresses.*.allocation_id", "aws_eip.test3", "id"),
				),
			},
		},
	})
}

func TestAccVPCNATGateway_perVPCLimit(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckNATGatewayDestroy,
		Steps: []resource.TestStep{
			{
				// Only one NAT gateway is allowed per VPC, so creating second in same VPC must fail
				Config:      testAccNATGatewayConfigPerVPCLimit(rName),
				ExpectError: regexp.MustCompile(`already exists`),
			},
		},
	})
}

func TestAccVPCNATGateway_tags(t *testing.T) {
	var natGateway ec2.NatGateway
	resourceName := "aws_nat_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckNATGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNATGatewayConfigTags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNATGatewayExists(resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNATGatewayConfigTags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNATGatewayExists(resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccNATGatewayConfigTags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNATGatewayExists(resourceName, &natGateway),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckNATGatewayDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_nat_gateway" {
			continue
		}

		_, err := tfec2.FindNATGatewayByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("EC2 NAT Gateway %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckNATGatewayExists(n string, v *ec2.NatGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 NAT Gateway ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		output, err := tfec2.FindNATGatewayByID(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckNATGatewayNotRecreated(before, after *ec2.NatGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.StringValue(before.NatGatewayId), aws.StringValue(after.NatGatewayId); before != after {
			return fmt.Errorf("EC2 NAT Gateway was recreated: %s -> %s", before, after)
		}

		return nil
	}
}

func testAccPreCheckAvailabilityZonesAtLeast(t *testing.T, count int) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	output, err := conn.DescribeAvailabilityZones(&ec2.DescribeAvailabilityZonesInput{
		Filters: tfec2.BuildAttributeFilterList(map[string]string{
			"opt-in-status": "opt-in-not-required",
			"state":         "available",
		}),
	})

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("error describing EC2 Availability Zones: %s", err)
	}

	if l := len(output.AvailabilityZones); l < count {
		t.Skipf("skipping acceptance testing: %d Availability Zone(s) available, %d required", l, count)
	}
}

func testAccNATGatewayConfigBase(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptIn(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

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
`, rName))
}

func testAccNATGatewayConfigBaseWithEIP(rName string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfigBase(rName), fmt.Sprintf(`
resource "aws_eip" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_eip" "test2" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_eip" "test3" {
  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAccNATGatewayConfig(rName string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfigBaseWithEIP(rName), `
resource "aws_nat_gateway" "test" {
  vpc_id = aws_vpc.test.id

  availability_zone_addresses {
    availability_zone = data.aws_availability_zones.available.names[0]
    allocation_id     = aws_eip.test.id
  }

  depends_on = [aws_internet_gateway.test]
}
`)
}

func testAccNATGatewayConfigAutoAllocate(rName string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfigBase(rName), fmt.Sprintf(`
resource "aws_nat_gateway" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }

  depends_on = [aws_internet_gateway.test]
}
`, rName))
}

func testAccNATGatewayConfigPerVPCLimit(rName string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfigBase(rName), `
resource "aws_nat_gateway" "test" {
  vpc_id = aws_vpc.test.id

  depends_on = [aws_internet_gateway.test]
}

resource "aws_nat_gateway" "test2" {
  vpc_id = aws_vpc.test.id

  depends_on = [aws_nat_gateway.test]
}
`)
}

func testAccNATGatewayConfigAvailabilityZoneAddresses1(rName string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfigBaseWithEIP(rName), `
resource "aws_nat_gateway" "test" {
  vpc_id = aws_vpc.test.id

  availability_zone_addresses {
    availability_zone = data.aws_availability_zones.available.names[0]
    allocation_id     = aws_eip.test.id
  }

  depends_on = [aws_internet_gateway.test]
}
`)
}

func testAccNATGatewayConfigAvailabilityZoneAddresses2(rName string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfigBaseWithEIP(rName), `
resource "aws_nat_gateway" "test" {
  vpc_id = aws_vpc.test.id

  availability_zone_addresses {
    availability_zone = data.aws_availability_zones.available.names[0]
    allocation_id     = aws_eip.test.id
  }

  availability_zone_addresses {
    availability_zone = data.aws_availability_zones.available.names[1]
    allocation_id     = aws_eip.test2.id
  }

  depends_on = [aws_internet_gateway.test]
}
`)
}

func testAccNATGatewayConfigAvailabilityZoneAddresses3(rName string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfigBaseWithEIP(rName), `
resource "aws_nat_gateway" "test" {
  vpc_id = aws_vpc.test.id

  availability_zone_addresses {
    availability_zone = data.aws_availability_zones.available.names[0]
    allocation_id     = aws_eip.test3.id
  }

  availability_zone_addresses {
    availability_zone = data.aws_availability_zones.available.names[1]
    allocation_id     = aws_eip.test2.id
  }

  depends_on = [aws_internet_gateway.test]
}
`)
}

func testAccNATGatewayConfigTags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfigBaseWithEIP(rName), fmt.Sprintf(`
resource "aws_nat_gateway" "test" {
  vpc_id = aws_vpc.test.id

  availability_zone_addresses {
    availability_zone = data.aws_availability_zones.available.names[0]
    allocation_id     = aws_eip.test.id
  }

  tags = {
    %[1]q = %[2]q
  }

  depends_on = [aws_internet_gateway.test]
}
`, tagKey1, tagValue1))
}

func testAccNATGatewayConfigTags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfigBaseWithEIP(rName), fmt.Sprintf(`
resource "aws_nat_gateway" "test" {
  vpc_id = aws_vpc.test.id

  availability_zone_addresses {
    availability_zone = data.aws_availability_zones.available.names[0]
    allocation_id     = aws_eip.test.id
  }

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }

  depends_on = [aws_internet_gateway.test]
}
`, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccNATGatewayConfigRoute(rName, destinationCidr string) string {
	return acctest.ConfigCompose(testAccNATGatewayConfig(rName), fmt.Sprintf(`
resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route" "test" {
  route_table_id         = aws_route_table.test.id
  destination_cidr_block = %[2]q
  nat_gateway_id         = aws_nat_gateway.test.id
}
`, rName, destinationCidr))
}
