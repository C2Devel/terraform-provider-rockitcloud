package ec2_test

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// This will currently skip EIPs with associations,
// although we depend on aws_vpc to potentially have
// the majority of those associations removed.

func TestAccEC2EIP_basic(t *testing.T) {
	var conf ec2.Address
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					// Not attached anywhere, so no DNS name to report.
					resource.TestCheckResourceAttr(resourceName, "public_dns", ""),
					resource.TestCheckResourceAttr(resourceName, "domain", ec2.DomainTypeVpc),
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

func TestAccEC2EIP_disappears(t *testing.T) {
	var conf ec2.Address
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceEIP(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccEC2EIP_instance(t *testing.T) {
	var conf ec2.Address
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPInstanceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					// Attached, so DNS names must be real, not empty.
					testAccCheckEIPDNS(resourceName, &conf),
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

// Regression test for https://github.com/hashicorp/terraform/issues/3429 (now
// https://github.com/hashicorp/terraform-provider-aws/issues/42)
func TestAccEC2EIP_Instance_reassociate(t *testing.T) {
	instanceResourceName := "aws_instance.test"
	resourceName := "aws_eip.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPInstanceReassociateConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "instance", instanceResourceName, "id"),
				),
			},
			{
				Config: testAccEIPInstanceReassociateConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "instance", instanceResourceName, "id"),
				),
				Taint: []string{resourceName},
			},
		},
	})
}

// Moving an EIP from one instance to another must keep it associated.
func TestAccEC2EIP_Instance_switch(t *testing.T) {
	var one ec2.Address
	resourceName := "aws_eip.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPInstanceAssociatedConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &one),
					testAccCheckEIPAttributes(&one),
					testAccCheckEIPAssociated(&one),
					resource.TestCheckResourceAttr(resourceName, "domain", ec2.DomainTypeVpc),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEIPInstanceAssociatedSwitchConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &one),
					testAccCheckEIPAttributes(&one),
					testAccCheckEIPAssociated(&one),
					resource.TestCheckResourceAttr(resourceName, "domain", ec2.DomainTypeVpc),
				),
			},
		},
	})
}

func TestAccEC2EIP_Instance_notAssociated(t *testing.T) {
	var conf ec2.Address
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPInstanceAssociateNotAssociatedConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEIPInstanceAssociateAssociatedConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					testAccCheckEIPAssociated(&conf),
				),
			},
		},
	})
}

func TestAccEC2EIP_networkInterface(t *testing.T) {
	var conf ec2.Address
	resourceName := "aws_eip.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPNetworkInterfaceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					testAccCheckEIPAssociated(&conf),
					// DNS names come from the attached interface, so they are
					// reported even when it belongs to no instance.
					resource.TestCheckResourceAttrPair(resourceName, "private_dns", "aws_network_interface.test", "private_dns_name"),
					resource.TestCheckResourceAttrSet(resourceName, "public_dns"),
					resource.TestCheckResourceAttrSet(resourceName, "allocation_id"),
					resource.TestCheckResourceAttr(resourceName, "domain", ec2.DomainTypeVpc),
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

// AssociateAddress ignores the requested private IP address and always uses the
// primary one, so the second EIP fails with Resource.AlreadyAssociated.
func TestAccEC2EIP_NetworkInterface_twoEIPsOneInterface(t *testing.T) {
	t.Skip("an EIP can only be associated with the primary private IP address of an interface")

	var one, two ec2.Address
	resourceName := "aws_eip.test"
	resourceName2 := "aws_eip.test2"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPMultiNetworkInterfaceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &one),
					testAccCheckEIPAttributes(&one),
					testAccCheckEIPAssociated(&one),
					resource.TestCheckResourceAttr(resourceName, "domain", ec2.DomainTypeVpc),

					testAccCheckEIPExists(resourceName2, &two),
					testAccCheckEIPAttributes(&two),
					testAccCheckEIPAssociated(&two),
					resource.TestCheckResourceAttr(resourceName2, "domain", ec2.DomainTypeVpc),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"associate_with_private_ip"},
			},
		},
	})
}

func TestAccEC2EIP_TagsEC2VPC_withVPCTrue(t *testing.T) {
	var conf ec2.Address
	resourceName := "aws_eip.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPTagsEC2VPCConfig(rName, "vpc = true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "domain", ec2.DomainTypeVpc),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.RandomName", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.TestName", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEIPTagsEC2VPCConfig(rName2, "vpc = true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.RandomName", rName2),
					resource.TestCheckResourceAttr(resourceName, "tags.TestName", rName2),
				),
			},
		},
	})
}

// Regression test for https://github.com/hashicorp/terraform-provider-aws/issues/18756
func TestAccEC2EIP_TagsEC2VPC_withoutVPCTrue(t *testing.T) {
	var conf ec2.Address
	resourceName := "aws_eip.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPTagsEC2VPCConfig(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "domain", ec2.DomainTypeVpc),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.RandomName", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.TestName", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEIPTagsEC2VPCConfig(rName2, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.RandomName", rName2),
					resource.TestCheckResourceAttr(resourceName, "tags.TestName", rName2),
				),
			},
		},
	})
}

func TestAccEC2EIP_PublicIPv4Pool_default(t *testing.T) {
	var conf ec2.Address
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPPublicIPv4PoolDefaultConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "public_ipv4_pool", "default"),
					resource.TestCheckResourceAttr(resourceName, "domain", ec2.DomainTypeVpc),
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

func TestAccEC2EIP_PublicIPv4Pool_custom(t *testing.T) {
	if os.Getenv("AWS_EC2_EIP_PUBLIC_IPV4_POOL") == "" {
		t.Skip("Environment variable AWS_EC2_EIP_PUBLIC_IPV4_POOL is not set")
	}

	var conf ec2.Address
	resourceName := "aws_eip.test"

	poolName := os.Getenv("AWS_EC2_EIP_PUBLIC_IPV4_POOL")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPPublicIPv4PoolCustomConfig(poolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "public_ipv4_pool", poolName),
					resource.TestCheckResourceAttr(resourceName, "domain", ec2.DomainTypeVpc),
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

func TestAccEC2EIP_customerOwnedIPv4Pool(t *testing.T) {
	var conf ec2.Address
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t); acctest.PreCheckOutpostsOutposts(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPCustomerOwnedIPv4PoolConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					resource.TestMatchResourceAttr(resourceName, "customer_owned_ipv4_pool", regexp.MustCompile(`^ipv4pool-coip-.+$`)),
					resource.TestMatchResourceAttr(resourceName, "customer_owned_ip", regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)),
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

// The API ignores network_border_group in AllocateAddress/ReleaseAddress
// and never returns it in DescribeAddresses.
func TestAccEC2EIP_networkBorderGroup(t *testing.T) {
	t.Skip("network_border_group is not supported")

	var conf ec2.Address
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPNetworkBorderGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "public_ipv4_pool", "default"),
					resource.TestCheckResourceAttr(resourceName, "network_border_group", acctest.Region()),
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

// testAccPreCheckWavelengthZoneAvailable's zone-type filter is not honored
// here, so the test does not skip on its own.
func TestAccEC2EIP_carrierIP(t *testing.T) {
	t.Skip("carrier_ip requires AWS Wavelength zones, which are not supported")

	var conf ec2.Address
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t); testAccPreCheckWavelengthZoneAvailable(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPCarrierIPConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					resource.TestCheckResourceAttrSet(resourceName, "carrier_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "network_border_group"),
					resource.TestCheckResourceAttr(resourceName, "public_ip", ""),
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

func TestAccEC2EIP_BYOIPAddress_default(t *testing.T) {
	// Test case address not set
	var conf ec2.Address
	resourceName := "aws_eip.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPConfig_BYOIPAddress_custom_default,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
				),
			},
		},
	})
}

func TestAccEC2EIP_BYOIPAddress_custom(t *testing.T) {
	// Test Case for address being set

	if os.Getenv("AWS_EC2_EIP_BYOIP_ADDRESS") == "" {
		t.Skip("Environment variable AWS_EC2_EIP_BYOIP_ADDRESS is not set")
	}

	var conf ec2.Address
	resourceName := "aws_eip.test"

	address := os.Getenv("AWS_EC2_EIP_BYOIP_ADDRESS")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPConfig_BYOIPAddress_custom(address),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "public_ip", address),
				),
			},
		},
	})
}

func TestAccEC2EIP_BYOIPAddress_customWithPublicIPv4Pool(t *testing.T) {
	// Test Case for both address and public_ipv4_pool being set
	if os.Getenv("AWS_EC2_EIP_BYOIP_ADDRESS") == "" {
		t.Skip("Environment variable AWS_EC2_EIP_BYOIP_ADDRESS is not set")
	}

	if os.Getenv("AWS_EC2_EIP_PUBLIC_IPV4_POOL") == "" {
		t.Skip("Environment variable AWS_EC2_EIP_PUBLIC_IPV4_POOL is not set")
	}

	var conf ec2.Address
	resourceName := "aws_eip.test"

	address := os.Getenv("AWS_EC2_EIP_BYOIP_ADDRESS")
	poolName := os.Getenv("AWS_EC2_EIP_PUBLIC_IPV4_POOL")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCheckEIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEIPConfig_BYOIPAddress_custom_with_PublicIPv4Pool(address, poolName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEIPExists(resourceName, &conf),
					testAccCheckEIPAttributes(&conf),
					resource.TestCheckResourceAttr(resourceName, "public_ip", address),
					resource.TestCheckResourceAttr(resourceName, "public_ipv4_pool", poolName),
				),
			},
		},
	})
}

func testAccCheckEIPDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_eip" {
			continue
		}

		if strings.Contains(rs.Primary.ID, "eipalloc") {
			req := &ec2.DescribeAddressesInput{
				AllocationIds: []*string{aws.String(rs.Primary.ID)},
			}
			describe, err := conn.DescribeAddresses(req)
			if err != nil {
				// Verify the error is what we want
				if tfawserr.ErrCodeEquals(err, tfec2.ErrCodeInvalidAllocationIDNotFound, tfec2.ErrCodeInvalidAddressNotFound) {
					continue
				}
				return err
			}

			if len(describe.Addresses) > 0 {
				return fmt.Errorf("still exists")
			}
		} else {
			req := &ec2.DescribeAddressesInput{
				PublicIps: []*string{aws.String(rs.Primary.ID)},
			}
			describe, err := conn.DescribeAddresses(req)
			if err != nil {
				// Verify the error is what we want
				if tfawserr.ErrCodeEquals(err, tfec2.ErrCodeInvalidAllocationIDNotFound, tfec2.ErrCodeInvalidAddressNotFound) {
					continue
				}
				return err
			}

			if len(describe.Addresses) > 0 {
				return fmt.Errorf("still exists")
			}
		}
	}

	return nil
}

func testAccCheckEIPAttributes(conf *ec2.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *conf.PublicIp == "" {
			return fmt.Errorf("empty public_ip")
		}

		return nil
	}
}

func testAccCheckEIPAssociated(conf *ec2.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if conf.AssociationId == nil || *conf.AssociationId == "" {
			return fmt.Errorf("empty association_id")
		}

		return nil
	}
}

// Fetches the network interface via the API, since its Terraform state isn't
// refreshed by the EIP's own apply and would show stale values.
func testAccCheckEIPDNS(resourceName string, conf *ec2.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(conf.NetworkInterfaceId) == "" {
			return fmt.Errorf("EIP is not attached to a network interface")
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		eni, err := tfec2.FindNetworkInterfaceByID(conn, aws.StringValue(conf.NetworkInterfaceId))
		if err != nil {
			return fmt.Errorf("error describing network interface %s: %w", aws.StringValue(conf.NetworkInterfaceId), err)
		}

		privateDNS := aws.StringValue(eni.PrivateDnsName)
		if privateDNS == "" {
			return fmt.Errorf("empty private_dns on network interface %s", aws.StringValue(conf.NetworkInterfaceId))
		}
		if got := rs.Primary.Attributes["private_dns"]; got != privateDNS {
			return fmt.Errorf("EIP private_dns %q does not match network interface private_dns %q", got, privateDNS)
		}

		if eni.Association == nil {
			return fmt.Errorf("no association on network interface %s", aws.StringValue(conf.NetworkInterfaceId))
		}

		publicDNS := aws.StringValue(eni.Association.PublicDnsName)
		if publicDNS == "" {
			return fmt.Errorf("empty public_dns on network interface %s", aws.StringValue(conf.NetworkInterfaceId))
		}
		if got := rs.Primary.Attributes["public_dns"]; got != publicDNS {
			return fmt.Errorf("EIP public_dns %q does not match network interface public_dns %q", got, publicDNS)
		}

		return nil
	}
}

func testAccCheckEIPExists(n string, res *ec2.Address) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EIP ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		input := &ec2.DescribeAddressesInput{}

		if strings.Contains(rs.Primary.ID, "eipalloc") {
			input.AllocationIds = aws.StringSlice([]string{rs.Primary.ID})
		} else {
			input.PublicIps = aws.StringSlice([]string{rs.Primary.ID})
		}

		var output *ec2.DescribeAddressesOutput

		err := resource.Retry(15*time.Minute, func() *resource.RetryError {
			var err error

			output, err = conn.DescribeAddresses(input)

			if tfawserr.ErrCodeEquals(err, tfec2.ErrCodeInvalidAllocationIDNotFound, tfec2.ErrCodeInvalidAddressNotFound) {
				return resource.RetryableError(err)
			}

			if err != nil {
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if tfresource.TimedOut(err) {
			output, err = conn.DescribeAddresses(input)
		}

		if err != nil {
			return fmt.Errorf("while describing addresses (%s): %w", rs.Primary.ID, err)
		}

		if len(output.Addresses) != 1 {
			return fmt.Errorf("wrong number of EIP found for (%s): %d", rs.Primary.ID, len(output.Addresses))
		}

		if aws.StringValue(output.Addresses[0].AllocationId) != rs.Primary.ID && aws.StringValue(output.Addresses[0].PublicIp) != rs.Primary.ID {
			return fmt.Errorf("EIP (%s) not found", rs.Primary.ID)
		}

		*res = *output.Addresses[0]

		return nil
	}
}

const testAccEIPConfig = `
resource "aws_eip" "test" {
}
`

func testAccEIPTagsEC2VPCConfig(rName, vpcConfig string) string {
	return fmt.Sprintf(`
resource "aws_eip" "test" {
  %[1]s

  tags = {
    RandomName = %[2]q
    TestName   = %[2]q
  }
}
`, vpcConfig, rName)
}

const testAccEIPPublicIPv4PoolDefaultConfig = `
resource "aws_eip" "test" {
}
`

func testAccEIPPublicIPv4PoolCustomConfig(poolName string) string {
	return fmt.Sprintf(`
resource "aws_eip" "test" {
  public_ipv4_pool = %[1]q
}
`, poolName)
}

const testAccEIPConfig_BYOIPAddress_custom_default = `
resource "aws_eip" "test" {
}
`

func testAccEIPConfig_BYOIPAddress_custom(address string) string {
	return fmt.Sprintf(`
resource "aws_eip" "test" {
  address = %[1]q
}
`, address)
}

func testAccEIPConfig_BYOIPAddress_custom_with_PublicIPv4Pool(address string, poolname string) string {
	return fmt.Sprintf(`
resource "aws_eip" "test" {
  address          = %[1]q
  public_ipv4_pool = %[2]q
}
`, address, poolname)
}

func testAccEIPInstanceAMIConfig() string {
	return `
data "aws_ami" "eip_test" {
  most_recent = true
  owners      = ["k2"]
}
`
}

func testAccEIPInstanceConfig() string {
	return acctest.ConfigCompose(
		testAccEIPInstanceAMIConfig(),
		acctest.ConfigAvailableAZsNoOptIn(),
		`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = cidrsubnet(aws_vpc.test.cidr_block, 8, 0)
  vpc_id            = aws_vpc.test.id
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.eip_test.id
  instance_type = "m1.micro"
  subnet_id     = aws_subnet.test.id
}

resource "aws_eip" "test" {
  instance = aws_instance.test.id
}
`)
}

func testAccEIPInstanceAssociatedConfig(rName string) string {
	return acctest.ConfigCompose(
		testAccEIPInstanceAMIConfig(),
		acctest.ConfigAvailableAZsNoOptIn(),
		fmt.Sprintf(`
resource "aws_vpc" "default" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.default.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  vpc_id                  = aws_vpc.default.id
  availability_zone       = data.aws_availability_zones.available.names[0]
  cidr_block              = "10.0.0.0/24"
  map_public_ip_on_launch = true

  depends_on = [aws_internet_gateway.gw]

  tags = {
    Name = %[1]q
  }
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.eip_test.id
  instance_type = "m1.micro"

  private_ip = "10.0.0.12"
  subnet_id  = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_instance" "test2" {
  ami           = data.aws_ami.eip_test.id
  instance_type = "m1.micro"

  private_ip = "10.0.0.19"
  subnet_id  = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_eip" "test" {
  instance = aws_instance.test2.id
}
`, rName))
}

func testAccEIPInstanceAssociatedSwitchConfig(rName string) string {
	return acctest.ConfigCompose(
		testAccEIPInstanceAMIConfig(),
		fmt.Sprintf(`
resource "aws_vpc" "default" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.default.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  vpc_id                  = aws_vpc.default.id
  cidr_block              = "10.0.0.0/24"
  map_public_ip_on_launch = true

  depends_on = [aws_internet_gateway.gw]

  tags = {
    Name = %[1]q
  }
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.eip_test.id
  instance_type = "m1.micro"

  private_ip = "10.0.0.12"
  subnet_id  = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_instance" "test2" {
  ami = data.aws_ami.eip_test.id

  instance_type = "m1.micro"

  private_ip = "10.0.0.19"
  subnet_id  = aws_subnet.test.id

  tags = {
    Name = "%[1]s-2"
  }
}

resource "aws_eip" "test" {
  instance = aws_instance.test.id
}
`, rName))
}

func testAccEIPNetworkInterfaceConfig(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigAvailableAZsNoOptIn(),
		fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/24"
  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id
}

resource "aws_subnet" "test" {
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  tags = {
    Name = %[1]q
  }
}

resource "aws_network_interface" "test" {
  subnet_id       = aws_subnet.test.id
  private_ips     = ["10.0.0.10"]
  security_groups = [aws_vpc.test.default_security_group_id]
}

resource "aws_eip" "test" {
  network_interface = aws_network_interface.test.id
  depends_on        = [aws_internet_gateway.test]
}
`, rName))
}

func testAccEIPMultiNetworkInterfaceConfig(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigAvailableAZsNoOptIn(),
		fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/24"
  tags = {
    Name = %[1]q
  }
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id
}

resource "aws_subnet" "test" {
  vpc_id            = aws_vpc.test.id
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  tags = {
    Name = %[1]q
  }
}

resource "aws_network_interface" "test" {
  subnet_id       = aws_subnet.test.id
  private_ips     = ["10.0.0.10", "10.0.0.11"]
  security_groups = [aws_vpc.test.default_security_group_id]
}

resource "aws_eip" "test" {
  network_interface         = aws_network_interface.test.id
  associate_with_private_ip = "10.0.0.10"
  depends_on                = [aws_internet_gateway.test]
}

resource "aws_eip" "test2" {
  network_interface         = aws_network_interface.test.id
  associate_with_private_ip = "10.0.0.11"
  depends_on                = [aws_internet_gateway.test]
}
`, rName))
}

func testAccEIPInstanceReassociateConfig(rName string) string {
	return acctest.ConfigCompose(
		testAccEIPInstanceAMIConfig(),
		fmt.Sprintf(`
resource "aws_eip" "test" {
  instance = aws_instance.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.eip_test.id
  instance_type = "m1.micro"
  subnet_id     = aws_subnet.test.id

  tags = {
    Name = %[1]q
  }

  lifecycle {
    create_before_destroy = true
  }
}

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

resource "aws_subnet" "test" {
  cidr_block = "10.0.0.0/24"
  vpc_id     = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table" "test" {
  vpc_id = aws_vpc.test.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.test.id
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_route_table_association" "test" {
  subnet_id      = aws_subnet.test.id
  route_table_id = aws_route_table.test.id
}
`, rName))
}

func testAccEIPInstanceAssociateNotAssociatedConfig() string {
	return acctest.ConfigCompose(
		acctest.ConfigAvailableAZsNoOptIn(),
		testAccEIPInstanceAMIConfig(), `
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = cidrsubnet(aws_vpc.test.cidr_block, 8, 0)
  vpc_id            = aws_vpc.test.id
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.eip_test.id
  instance_type = "m1.micro"
  subnet_id     = aws_subnet.test.id
}

resource "aws_eip" "test" {
}
`)
}

func testAccEIPInstanceAssociateAssociatedConfig() string {
	return acctest.ConfigCompose(
		acctest.ConfigAvailableAZsNoOptIn(),
		testAccEIPInstanceAMIConfig(), `
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = cidrsubnet(aws_vpc.test.cidr_block, 8, 0)
  vpc_id            = aws_vpc.test.id
}

resource "aws_internet_gateway" "test" {
  vpc_id = aws_vpc.test.id
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.eip_test.id
  instance_type = "m1.micro"
  subnet_id     = aws_subnet.test.id
}

resource "aws_eip" "test" {
  instance = aws_instance.test.id
}
`)
}

func testAccEIPCustomerOwnedIPv4PoolConfig() string {
	return `
data "aws_ec2_coip_pools" "test" {}

resource "aws_eip" "test" {
  customer_owned_ipv4_pool = tolist(data.aws_ec2_coip_pools.test.pool_ids)[0]
}
`
}

const testAccEIPNetworkBorderGroupConfig = `
data "aws_region" current {}

resource "aws_eip" "test" {
  network_border_group = data.aws_region.current.name
}
`

func testAccEIPCarrierIPConfig(rName string) string {
	return acctest.ConfigCompose(
		testAccAvailableAZsWavelengthZonesDefaultExcludeConfig(),
		fmt.Sprintf(`
data "aws_availability_zone" "available" {
  name = data.aws_availability_zones.available.names[0]
}

resource "aws_eip" "test" {
  network_border_group = data.aws_availability_zone.available.network_border_group

  tags = {
    Name = %[1]q
  }
}
`, rName))
}
