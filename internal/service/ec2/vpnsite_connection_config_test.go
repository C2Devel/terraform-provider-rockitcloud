package ec2

import (
	"encoding/xml"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const testVpnConnectionCustomerGatewayConfiguration = `<?xml version='1.0' encoding='UTF-8'?>
<vpn_connection id="vpn-12345678">
  <customer_gateway_id>cgw-12345678</customer_gateway_id>
  <vpn_gateway_id>vgw-12345678</vpn_gateway_id>
  <vpn_connection_type>ipsec.1</vpn_connection_type>
  <ipv4_network_cidr_local>10.1.0.0/16</ipv4_network_cidr_local>
  <ipv4_network_cidr_remote>10.2.0.0/16</ipv4_network_cidr_remote>
  <ipsec_tunnel>
    <customer_gateway>
      <tunnel_outside_address><ip_address>198.51.100.10</ip_address></tunnel_outside_address>
      <tunnel_inside_address><ip_address>169.254.252.2</ip_address></tunnel_inside_address>
      <bgp><asn>65000</asn><hold_time>30</hold_time></bgp>
    </customer_gateway>
    <vpn_gateway>
      <tunnel_outside_address><ip_address>203.0.113.1</ip_address></tunnel_outside_address>
      <tunnel_inside_address><ip_address>169.254.252.1</ip_address></tunnel_inside_address>
      <bgp><asn>64512</asn><hold_time>30</hold_time></bgp>
    </vpn_gateway>
    <ike><pre_shared_key>FIRST_KEY</pre_shared_key></ike>
    <tunnel_cidr>169.254.252.0/30</tunnel_cidr>
  </ipsec_tunnel>
  <ipsec_tunnel>
    <customer_gateway>
      <tunnel_outside_address><ip_address>198.51.100.10</ip_address></tunnel_outside_address>
      <tunnel_inside_address><ip_address>169.254.252.6</ip_address></tunnel_inside_address>
      <bgp><asn>65000</asn><hold_time>30</hold_time></bgp>
    </customer_gateway>
    <vpn_gateway>
      <tunnel_outside_address><ip_address>203.0.113.2</ip_address></tunnel_outside_address>
      <tunnel_inside_address><ip_address>169.254.252.5</ip_address></tunnel_inside_address>
      <bgp><asn>64512</asn><hold_time>30</hold_time></bgp>
    </vpn_gateway>
    <ike><pre_shared_key>SECOND_KEY</pre_shared_key></ike>
    <tunnel_cidr>169.254.252.4/30</tunnel_cidr>
  </ipsec_tunnel>
</vpn_connection>`

func TestXmlVpnConnectionConfig(t *testing.T) {
	var config XmlVpnConnectionConfig

	if err := xml.Unmarshal([]byte(testVpnConnectionCustomerGatewayConfiguration), &config); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if got, want := len(config.Tunnels), 2; got != want {
		t.Fatalf("len(Tunnels) = %d, want %d", got, want)
	}

	want := XmlIpsecTunnel{
		BGPASN:           "64512",
		BGPHoldTime:      30,
		CgwInsideAddress: "169.254.252.6",
		OutsideAddress:   "203.0.113.2",
		VgwInsideAddress: "169.254.252.5",
	}

	if got := config.Tunnels[1]; got != want {
		t.Errorf("Tunnels[1] = %+v, want %+v", got, want)
	}
}

func TestXmlVpnConnectionConfig_tunnel(t *testing.T) {
	var config XmlVpnConnectionConfig

	if err := xml.Unmarshal([]byte(testVpnConnectionCustomerGatewayConfiguration), &config); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	testCases := []struct {
		Name             string
		Index            int
		TunnelOption     *ec2.TunnelOption
		VgwInsideAddress string
	}{
		{
			Name:             "matched by outside address",
			Index:            0,
			TunnelOption:     &ec2.TunnelOption{OutsideIpAddress: aws.String("203.0.113.2")},
			VgwInsideAddress: "169.254.252.5",
		},
		{
			Name:             "unknown outside address",
			Index:            1,
			TunnelOption:     &ec2.TunnelOption{OutsideIpAddress: aws.String("203.0.113.3")},
			VgwInsideAddress: "169.254.252.5",
		},
		{
			Name:             "no outside address",
			Index:            0,
			TunnelOption:     &ec2.TunnelOption{},
			VgwInsideAddress: "169.254.252.1",
		},
		{
			Name:             "no tunnel option",
			Index:            1,
			VgwInsideAddress: "169.254.252.5",
		},
		{
			Name:  "no tunnel",
			Index: 2,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			tunnel := config.tunnel(testCase.Index, testCase.TunnelOption)

			if testCase.VgwInsideAddress == "" {
				if tunnel != nil {
					t.Fatalf("got %+v, want nil", tunnel)
				}

				return
			}

			if tunnel == nil {
				t.Fatal("got nil")
			}

			if got, want := tunnel.VgwInsideAddress, testCase.VgwInsideAddress; got != want {
				t.Errorf("VgwInsideAddress = %q, want %q", got, want)
			}
		})
	}
}

func TestValidateVpnConnectionTunnelInsideCIDR(t *testing.T) {
	validCidrs := []string{
		"169.254.252.0/30",
		"169.254.252.2/30",
		"169.254.255.252/30",
	}
	for _, v := range validCidrs {
		_, errors := validateVpnConnectionTunnelInsideCIDR(v, "tunnel1_inside_cidr")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid tunnel inside CIDR: %q", v, errors)
		}
	}

	invalidCidrs := []string{
		"",
		"169.254.252.0",
		"169.254.252.0/29",
		"169.254.252.0/31",
		"169.254.251.252/30",
		"169.254.10.0/30",
		"10.0.0.0/30",
		"fd00::/126",
	}
	for _, v := range invalidCidrs {
		_, errors := validateVpnConnectionTunnelInsideCIDR(v, "tunnel1_inside_cidr")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid tunnel inside CIDR", v)
		}
	}
}
