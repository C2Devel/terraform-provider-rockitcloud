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

func TestAccVPNSiteConnection_basic(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionConfig(rName, rBgpAsn),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "ec2", regexp.MustCompile(`vpn-connection/vpn-.+`)),
					resource.TestCheckResourceAttrSet(resourceName, "customer_gateway_configuration"),
					resource.TestCheckResourceAttrPair(resourceName, "customer_gateway_id", "aws_customer_gateway.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "local_ipv4_network_cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "remote_ipv4_network_cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "static_routes_only", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tunnel_inside_ip_version", "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_address"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_bgp_asn"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_bgp_holdtime", "30"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_cgw_inside_address"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_ike_versions.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_ike_versions.*", "ikev1"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_inside_cidr"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase1_dh_group_numbers.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase1_encryption_algorithms.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase1_integrity_algorithms.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase1_lifetime_seconds", "28800"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase2_dh_group_numbers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase2_dh_group_numbers.*", "14"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase2_encryption_algorithms.#", "10"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase2_integrity_algorithms.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase2_lifetime_seconds", "3600"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_preshared_key"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_replay_window_size", "1024"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_vgw_inside_address"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_address"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_bgp_asn"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_bgp_holdtime", "30"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_cgw_inside_address"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_ike_versions.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_inside_cidr"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_phase1_lifetime_seconds", "28800"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_phase2_lifetime_seconds", "3600"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_preshared_key"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_replay_window_size", "1024"),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_vgw_inside_address"),
					resource.TestCheckResourceAttr(resourceName, "type", "ipsec.1"),
					resource.TestCheckResourceAttr(resourceName, "vgw_telemetry.#", "2"),
					testAccCheckVPNConnectionVPNGatewayID(resourceName, "aws_vpc.test"),
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

func TestAccVPNSiteConnection_disappears(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionConfig(rName, rBgpAsn),
				Check: resource.ComposeTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn),
					acctest.CheckResourceDisappears(acctest.Provider, tfec2.ResourceVPNConnection(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVPNSiteConnection_tags(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn1, vpn2 ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionTags1Config(rName, rBgpAsn, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn1),
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
				Config: testAccVPNConnectionTags2Config(rName, rBgpAsn, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn2),
					testAccCheckVPNConnectionNotRecreated(&vpn1, &vpn2),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccVPNConnectionTags1Config(rName, rBgpAsn, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn2),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccVPNSiteConnection_tunnel2Options(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccVPNConnectionTunnel2OptionsConfig(rName, rBgpAsn),
				ExpectError: regexp.MustCompile(`tunnel2_replay_window_size requires at least one tunnel1_\* argument`),
			},
		},
	})
}

func TestAccVPNSiteConnection_tunnelOptions(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn1, vpn2 ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionTunnelOptionsConfig(rName, rBgpAsn, 512),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn1),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_ike_versions.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_ike_versions.*", "ikev2"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_inside_cidr", "169.254.253.0/30"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase1_dh_group_numbers.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase1_dh_group_numbers.*", "19"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase1_dh_group_numbers.*", "20"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase1_encryption_algorithms.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase1_encryption_algorithms.*", "aes_gcm256"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase1_encryption_algorithms.*", "chacha20poly1305"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase1_integrity_algorithms.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase1_integrity_algorithms.*", "sha512"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase1_lifetime_seconds", "14400"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase2_dh_group_numbers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase2_dh_group_numbers.*", "21"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase2_encryption_algorithms.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase2_encryption_algorithms.*", "aes_gcm256"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase2_integrity_algorithms.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase2_integrity_algorithms.*", "sha384"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel1_phase2_integrity_algorithms.*", "sha512"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_phase2_lifetime_seconds", "1800"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_preshared_key", "Tunnel1.Preshared_Key"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_replay_window_size", "512"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel2_ike_versions.*", "ikev1"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_inside_cidr", "169.254.253.4/30"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_phase1_dh_group_numbers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel2_phase1_dh_group_numbers.*", "14"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_phase1_encryption_algorithms.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel2_phase1_encryption_algorithms.*", "aes128"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel2_phase1_encryption_algorithms.*", "aes256"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_phase1_lifetime_seconds", "900"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_phase2_dh_group_numbers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tunnel2_phase2_dh_group_numbers.*", "0"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_phase2_lifetime_seconds", "900"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_preshared_key", "Tunnel2.Preshared_Key"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_replay_window_size", "32"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// TODO: change modify
			{
				Config: testAccVPNConnectionTunnelOptionsConfig(rName, rBgpAsn, 2048),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn2),
					testAccCheckVPNConnectionRecreated(&vpn1, &vpn2),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_replay_window_size", "2048"),
				),
			},
		},
	})
}

func TestAccVPNSiteConnection_tunnel1Options(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn1, vpn2, vpn3 ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionTunnel1OptionsConfig(rName, rBgpAsn, "169.254.254.2/30", "Tunnel1.Preshared_Key"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn1),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel1_address"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_cgw_inside_address", "169.254.254.1"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_inside_cidr", "169.254.254.2/30"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_preshared_key", "Tunnel1.Preshared_Key"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_replay_window_size", "1024"),
					resource.TestCheckResourceAttr(resourceName, "tunnel1_vgw_inside_address", "169.254.254.2"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_address", ""),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_ike_versions.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_inside_cidr", ""),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_preshared_key", ""),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_replay_window_size", "0"),
					resource.TestCheckResourceAttr(resourceName, "vgw_telemetry.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccVPNConnectionTunnel1And2OptionsConfig(rName, rBgpAsn, "169.254.254.2/30", "Tunnel1.Preshared_Key", "Tunnel2.Preshared_Key"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn2),
					testAccCheckVPNConnectionRecreated(&vpn1, &vpn2),
					resource.TestCheckResourceAttrSet(resourceName, "tunnel2_address"),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_preshared_key", "Tunnel2.Preshared_Key"),
					resource.TestCheckResourceAttr(resourceName, "vgw_telemetry.#", "2"),
				),
			},
			{
				Config: testAccVPNConnectionTunnel1OptionsConfig(rName, rBgpAsn, "169.254.254.2/30", "Tunnel1.Preshared_Key"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn3),
					testAccCheckVPNConnectionRecreated(&vpn2, &vpn3),
					resource.TestCheckResourceAttr(resourceName, "tunnel2_address", ""),
					resource.TestCheckResourceAttr(resourceName, "vgw_telemetry.#", "1"),
				),
			},
		},
	})
}

func TestAccVPNSiteConnection_localRemoteIPv4CIDRs(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn1, vpn2 ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionLocalRemoteIPv4CIDRsConfig(rName, rBgpAsn, "10.111.0.0/16", "10.222.33.0/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn1),
					resource.TestCheckResourceAttr(resourceName, "local_ipv4_network_cidr", "10.111.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "remote_ipv4_network_cidr", "10.222.33.0/24"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// TODO: change modify
			{
				Config: testAccVPNConnectionLocalRemoteIPv4CIDRsConfig(rName, rBgpAsn, "10.112.0.0/16", "10.222.32.0/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn2),
					testAccCheckVPNConnectionRecreated(&vpn1, &vpn2),
					resource.TestCheckResourceAttr(resourceName, "local_ipv4_network_cidr", "10.112.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "remote_ipv4_network_cidr", "10.222.32.0/24"),
				),
			},
		},
	})
}

func TestAccVPNSiteConnection_vpnGatewayIDForm(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)
	resourceName := "aws_vpn_connection.test"
	var vpn1, vpn2 ec2.VpnConnection

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPNConnectionVPNGatewayIDFormConfig(rName, rBgpAsn, `replace(aws_vpc.test.id, "vpc-", "vgw-")`),
				Check: resource.ComposeTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn1),
					testAccCheckVPNConnectionVPNGatewayID(resourceName, "aws_vpc.test"),
				),
			},
			{
				Config: testAccVPNConnectionVPNGatewayIDFormConfig(rName, rBgpAsn, `aws_vpc.test.id`),
				Check: resource.ComposeTestCheckFunc(
					testAccVPNConnectionExists(resourceName, &vpn2),
					testAccCheckVPNConnectionNotRecreated(&vpn1, &vpn2),
					testAccCheckVPNConnectionVPNGatewayID(resourceName, "aws_vpc.test"),
				),
			},
		},
	})
}

func TestAccVPNSiteConnection_duplicate(t *testing.T) {
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rBgpAsn := sdkacctest.RandIntRange(64512, 65534)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		ErrorCheck:        acctest.ErrorCheck(t, ec2.EndpointsID),
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccVPNConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccVPNConnectionDuplicateConfig(rName, rBgpAsn),
				ExpectError: regexp.MustCompile(`EC2 VPN Connection \(vpn-.+\) between Customer Gateway \(cgw-.+\) and VPN Gateway \(vgw-.+\) already exists`),
			},
		},
	})
}

func testAccVPNConnectionDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_vpn_connection" {
			continue
		}

		_, err := tfec2.FindVPNConnectionByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("EC2 VPN Connection %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccVPNConnectionExists(n string, v *ec2.VpnConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 VPN Connection ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn

		output, err := tfec2.FindVPNConnectionByID(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckVPNConnectionNotRecreated(before, after *ec2.VpnConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.StringValue(before.VpnConnectionId), aws.StringValue(after.VpnConnectionId); before != after {
			return fmt.Errorf("Expected EC2 VPN Connection IDs not to change, but got before: %s, after: %s", before, after)
		}

		return nil
	}
}

func testAccCheckVPNConnectionRecreated(before, after *ec2.VpnConnection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.StringValue(before.VpnConnectionId), aws.StringValue(after.VpnConnectionId); before == after {
			return fmt.Errorf("Expected EC2 VPN Connection ID to change, but got: %s", before)
		}

		return nil
	}
}

func testAccCheckVPNConnectionVPNGatewayID(n, vpcResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpc, ok := s.RootModule().Resources[vpcResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", vpcResourceName)
		}

		if got, want := rs.Primary.Attributes["vpn_gateway_id"], tfec2.NormalizeVPNGatewayID(vpc.Primary.ID); got != want {
			return fmt.Errorf("vpn_gateway_id = %s, want %s", got, want)
		}

		return nil
	}
}

// VPN gateway is created together with VPC.
func testAccVPNConnectionBaseConfig(rName string, rBgpAsn int) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_customer_gateway" "test" {
  bgp_asn    = %[2]d
  ip_address = "178.0.0.1"
  type       = "ipsec.1"

  tags = {
    Name = %[1]q
  }
}
`, rName, rBgpAsn)
}

func testAccVPNConnectionConfig(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), `
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = aws_vpc.test.id
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"
}
`)
}

func testAccVPNConnectionTags1Config(rName string, rBgpAsn int, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), fmt.Sprintf(`
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = aws_vpc.test.id
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"

  tags = {
    %[1]q = %[2]q
  }
}
`, tagKey1, tagValue1))
}

func testAccVPNConnectionTags2Config(rName string, rBgpAsn int, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), fmt.Sprintf(`
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = aws_vpc.test.id
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }
}
`, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccVPNConnectionTunnel2OptionsConfig(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), `
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = aws_vpc.test.id
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"

  tunnel2_replay_window_size = 512
}
`)
}

func testAccVPNConnectionTunnelOptionsConfig(rName string, rBgpAsn int, tunnel1ReplayWindowSize int) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), fmt.Sprintf(`
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = aws_vpc.test.id
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"

  tunnel1_ike_versions                 = ["ikev2"]
  tunnel1_inside_cidr                  = "169.254.253.0/30"
  tunnel1_phase1_dh_group_numbers      = [19, 20]
  tunnel1_phase1_encryption_algorithms = ["aes_gcm256", "chacha20poly1305"]
  tunnel1_phase1_integrity_algorithms  = ["sha512"]
  tunnel1_phase1_lifetime_seconds      = 14400
  tunnel1_phase2_dh_group_numbers      = [21]
  tunnel1_phase2_encryption_algorithms = ["aes_gcm256"]
  tunnel1_phase2_integrity_algorithms  = ["sha384", "sha512"]
  tunnel1_phase2_lifetime_seconds      = 1800
  tunnel1_preshared_key                = "Tunnel1.Preshared_Key"
  tunnel1_replay_window_size           = %[1]d

  tunnel2_ike_versions                 = ["ikev1"]
  tunnel2_inside_cidr                  = "169.254.253.4/30"
  tunnel2_phase1_dh_group_numbers      = [14]
  tunnel2_phase1_encryption_algorithms = ["aes128", "aes256"]
  tunnel2_phase1_lifetime_seconds      = 900
  tunnel2_phase2_dh_group_numbers      = [0]
  tunnel2_phase2_lifetime_seconds      = 900
  tunnel2_preshared_key                = "Tunnel2.Preshared_Key"
  tunnel2_replay_window_size           = 32
}
`, tunnel1ReplayWindowSize))
}

func testAccVPNConnectionTunnel1OptionsConfig(rName string, rBgpAsn int, tunnel1InsideCidr, tunnel1PresharedKey string) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), fmt.Sprintf(`
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = aws_vpc.test.id
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"

  tunnel1_inside_cidr   = %[1]q
  tunnel1_preshared_key = %[2]q
}
`, tunnel1InsideCidr, tunnel1PresharedKey))
}

func testAccVPNConnectionTunnel1And2OptionsConfig(rName string, rBgpAsn int, tunnel1InsideCidr, tunnel1PresharedKey, tunnel2PresharedKey string) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), fmt.Sprintf(`
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = aws_vpc.test.id
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"

  tunnel1_inside_cidr   = %[1]q
  tunnel1_preshared_key = %[2]q

  tunnel2_preshared_key = %[3]q
}
`, tunnel1InsideCidr, tunnel1PresharedKey, tunnel2PresharedKey))
}

func testAccVPNConnectionLocalRemoteIPv4CIDRsConfig(rName string, rBgpAsn int, localIpv4Cidr string, remoteIpv4Cidr string) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), fmt.Sprintf(`
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = aws_vpc.test.id
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"

  local_ipv4_network_cidr  = %[1]q
  remote_ipv4_network_cidr = %[2]q
}
`, localIpv4Cidr, remoteIpv4Cidr))
}

func testAccVPNConnectionVPNGatewayIDFormConfig(rName string, rBgpAsn int, vpnGatewayID string) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), fmt.Sprintf(`
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = %[1]s
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"
}
`, vpnGatewayID))
}

func testAccVPNConnectionDuplicateConfig(rName string, rBgpAsn int) string {
	return acctest.ConfigCompose(testAccVPNConnectionBaseConfig(rName, rBgpAsn), `
resource "aws_vpn_connection" "test" {
  vpn_gateway_id      = aws_vpc.test.id
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"
}

resource "aws_vpn_connection" "duplicate" {
  vpn_gateway_id      = aws_vpn_connection.test.vpn_gateway_id
  customer_gateway_id = aws_customer_gateway.test.id
  type                = "ipsec.1"
}
`)
}
