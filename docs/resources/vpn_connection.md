---
subcategory: "VPN (Site-to-Site)"
layout: "aws"
page_title: "aws_vpn_connection"
description: |-
  Manages a Site-to-Site VPN connection.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[sensitive-data]: https://developer.hashicorp.com/terraform/language/state/sensitive-data
[vpn-connections]: https://docs.k2.cloud/en/services/interconnect/vpn_connections/operations.html
[create-vpn-connection]: https://docs.k2.cloud/en/api/ec2/actions/vpn_connections/CreateVpnConnection.html

# Resource: aws_vpn_connection

Manages a Site-to-Site VPN connection.
A Site-to-Site VPN connection is an IPsec VPN connection between a VPC and an on-premises network.

For more information about VPN connections, see [user documentation][vpn-connections].

~> **Note** All arguments including `tunnel1_preshared_key` and `tunnel2_preshared_key` are stored in the raw state as plain text.
[Read more about sensitive data in state][sensitive-data].

-> **Note** For convenience, the ID of the VPN gateway is the same as the ID of the VPC to which it belongs (`vpc-ABCD1234`/`vgw-ABCD1234`).
Either of them can be specified in the `vpn_gateway_id` argument.

~> **Note** Only one VPN connection can exist between a customer gateway and a VPC (VPN gateway).
An attempt to create another one succeeds, but the platform returns the existing VPN connection.
Do not manage the same pair of a customer gateway and a VPC with several `aws_vpn_connection` resources.

~> **Note** The provider supports only the high-availability VPN connections with two tunnels.
Specify the options for both tunnels or for none of them.
If the options are specified for one tunnel only, the platform creates a VPN connection with a single tunnel, and the provider fails to read it.

~> **Note** The platform does not support modifying a VPN connection except for its tags.
A change of any other argument of an existing VPN connection ends with the `Unsupported action` error and leaves the new value in the Terraform state.
To change the arguments, recreate the resource.

~> **Note** The tunnel options of a high-availability VPN connection are assigned to the tunnels in the order of the tunnel numbers.
For dynamic routing, enable the route propagation with the [`aws_vpn_gateway_route_propagation`](vpn_gateway_route_propagation.md) resource.

## Example usage

### Basic example

```terraform
resource "aws_vpc" "example" {
  cidr_block = "172.16.8.0/24"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_customer_gateway" "example" {
  bgp_asn    = 65000
  ip_address = "172.0.0.1"
  type       = "ipsec.1"

  tags = {
    Name = "tf-customer-gateway"
  }
}

resource "aws_vpn_connection" "example" {
  vpn_gateway_id      = aws_vpc.example.id
  customer_gateway_id = aws_customer_gateway.example.id
  type                = aws_customer_gateway.example.type

  tags = {
    Name = "tf-vpn-connection"
  }
}
```

### Specific example: tunnel options

```terraform
resource "aws_vpc" "example" {
  cidr_block = "172.16.8.0/24"
}

resource "aws_customer_gateway" "example" {
  bgp_asn    = 65000
  ip_address = "172.0.0.1"
  type       = "ipsec.1"
}

resource "aws_vpn_connection" "example" {
  vpn_gateway_id      = aws_vpc.example.id
  customer_gateway_id = aws_customer_gateway.example.id
  type                = aws_customer_gateway.example.type

  tunnel1_ike_versions            = ["ikev2"]
  tunnel1_phase1_lifetime_seconds = 14400
  tunnel1_replay_window_size      = 512

  tunnel2_ike_versions            = ["ikev2"]
  tunnel2_phase1_lifetime_seconds = 14400
  tunnel2_replay_window_size      = 512
}
```

## Argument reference

The following arguments are required:

* `customer_gateway_id` - (Required, String) The ID of the customer gateway.
* `type` - (Required, Forces new resource, String) The type of the VPN connection.
    * _Valid values:_ `ipsec.1`, `ipsec.legacy`
    * _Constraints:_ Must match the type of the customer gateway
* `vpn_gateway_id` - (Required, String) The ID of the VPN gateway or of the VPC to which it belongs.

The following arguments are optional:

* `local_ipv4_network_cidr` - (Optional, String) The IPv4 CIDR on the customer gateway (on-premises) side of the VPN connection.
    * _Default value:_ `0.0.0.0/0`
    * _Constraints:_ Must not fall within the 169.254.0.0/16 range
* `remote_ipv4_network_cidr` - (Optional, String) The IPv4 CIDR on the cloud side of the VPN connection.
    * _Default value:_ `0.0.0.0/0`
    * _Constraints:_ Must not fall within the 169.254.0.0/16 range
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource.
  If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

The tunnel options are described in the [API documentation][create-vpn-connection].
The following arguments configure the first VPN tunnel:

* `tunnel1_ike_versions` - (Optional, Set of strings) The IKE version permitted for the first VPN tunnel.
    * _Valid values:_ `ikev1`, `ikev2`
    * _Default value:_ `ikev1`
    * _Constraints:_ Only one IKE version can be specified, using both versions at once is not supported
* `tunnel1_inside_cidr` - (Optional, Forces new resource, String) The CIDR block of the inside IP addresses for the first VPN tunnel.
    * _Constraints:_ A /30 CIDR block from the 169.254.252.0/22 range, unique among all connections of the VPN gateway and of the public IP address of the customer gateway
    * _Default value:_ The block is assigned automatically
* `tunnel1_preshared_key` - (Optional, String) The preshared key (PSK) of the first VPN tunnel.
    * _Constraints:_ From 8 to 64 characters, alphanumeric characters, periods (.) and underscores (_) only, cannot start with zero (0)
    * _Default value:_ The key is generated automatically
* `tunnel1_phase1_dh_group_numbers` - (Optional, Set of integers) The Diffie-Hellman group numbers permitted for the first VPN tunnel for phase 1 of the IKE negotiation.
    * _Valid values:_ `2`, `5`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`
    * _Default value:_ `5`, `14`, `15`, `16`, `17`, `18`
    * _Constraints:_ No more than six values
* `tunnel1_phase1_encryption_algorithms` - (Optional, Set of strings) The encryption algorithms permitted for the first VPN tunnel for phase 1 of the IKE negotiation.
    * _Valid values:_ `aes128`, `aes256`, `aes_ctr128`, `aes_ctr256`, `aes_gcm128`, `aes_gcm256`, `camellia128`, `camellia256`, `chacha20poly1305`
    * _Default value:_ All algorithms supported by the selected IKE version
    * _Constraints:_ `aes_gcm128`, `aes_gcm256` and `chacha20poly1305` are not supported by `ikev1`
* `tunnel1_phase1_integrity_algorithms` - (Optional, Set of strings) The integrity algorithms permitted for the first VPN tunnel for phase 1 of the IKE negotiation.
    * _Valid values:_ `sha1`, `sha256`, `sha384`, `sha512`
    * _Default value:_ `sha1`, `sha256`, `sha384`, `sha512`
* `tunnel1_phase1_lifetime_seconds` - (Optional, Integer) The lifetime for phase 1 of the IKE negotiation for the first VPN tunnel, in seconds.
    * _Valid values:_ From `900` to `28800`
    * _Default value:_ `28800`
* `tunnel1_phase2_dh_group_numbers` - (Optional, Set of integers) The Diffie-Hellman group numbers permitted for the first VPN tunnel for phase 2 of the IKE negotiation.
    * _Valid values:_ `0`, `2`, `5`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`
    * _Default value:_ `14`
    * _Constraints:_ No more than one value, `19`, `20` and `21` are not supported by `ikev1`
* `tunnel1_phase2_encryption_algorithms` - (Optional, Set of strings) The encryption algorithms permitted for the first VPN tunnel for phase 2 of the IKE negotiation.
    * _Valid values:_ `aes128`, `aes256`, `aes_ccm128`, `aes_ccm256`, `aes_ctr128`, `aes_ctr256`, `aes_gcm128`, `aes_gcm256`, `camellia128`, `camellia256`, `chacha20poly1305`
    * _Default value:_ All algorithms supported by the selected IKE version
    * _Constraints:_ `chacha20poly1305` is not supported by `ikev1`
* `tunnel1_phase2_integrity_algorithms` - (Optional, Set of strings) The integrity algorithms permitted for the first VPN tunnel for phase 2 of the IKE negotiation.
    * _Valid values:_ `sha1`, `sha256`, `sha384`, `sha512`
    * _Default value:_ `sha1`, `sha256`, `sha384`, `sha512`
* `tunnel1_phase2_lifetime_seconds` - (Optional, Integer) The lifetime for phase 2 of the IKE negotiation for the first VPN tunnel, in seconds.
    * _Valid values:_ From `900` to `3600`
    * _Default value:_ `3600`
    * _Constraints:_ Cannot exceed the value of `tunnel1_phase1_lifetime_seconds`
* `tunnel1_replay_window_size` - (Optional, Integer) The number of packets in an IKE replay window for the first VPN tunnel.
    * _Valid values:_ From `32` to `2048`
    * _Default value:_ `1024`

The following arguments configure the second VPN tunnel:

* `tunnel2_ike_versions` - (Optional, Set of strings) The IKE version permitted for the second VPN tunnel.
    * _Valid values:_ `ikev1`, `ikev2`
    * _Default value:_ `ikev1`
    * _Constraints:_ Only one IKE version can be specified, using both versions at once is not supported
* `tunnel2_inside_cidr` - (Optional, Forces new resource, String) The CIDR block of the inside IP addresses for the second VPN tunnel.
    * _Constraints:_ A /30 CIDR block from the 169.254.252.0/22 range, unique among all connections of the VPN gateway and of the public IP address of the customer gateway
    * _Default value:_ The block is assigned automatically
* `tunnel2_preshared_key` - (Optional, String) The preshared key (PSK) of the second VPN tunnel.
    * _Constraints:_ From 8 to 64 characters, alphanumeric characters, periods (.) and underscores (_) only, cannot start with zero (0)
    * _Default value:_ The key is generated automatically
* `tunnel2_phase1_dh_group_numbers` - (Optional, Set of integers) The Diffie-Hellman group numbers permitted for the second VPN tunnel for phase 1 of the IKE negotiation.
    * _Valid values:_ `2`, `5`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`
    * _Default value:_ `5`, `14`, `15`, `16`, `17`, `18`
    * _Constraints:_ No more than six values
* `tunnel2_phase1_encryption_algorithms` - (Optional, Set of strings) The encryption algorithms permitted for the second VPN tunnel for phase 1 of the IKE negotiation.
    * _Valid values:_ `aes128`, `aes256`, `aes_ctr128`, `aes_ctr256`, `aes_gcm128`, `aes_gcm256`, `camellia128`, `camellia256`, `chacha20poly1305`
    * _Default value:_ All algorithms supported by the selected IKE version
    * _Constraints:_ `aes_gcm128`, `aes_gcm256` and `chacha20poly1305` are not supported by `ikev1`
* `tunnel2_phase1_integrity_algorithms` - (Optional, Set of strings) The integrity algorithms permitted for the second VPN tunnel for phase 1 of the IKE negotiation.
    * _Valid values:_ `sha1`, `sha256`, `sha384`, `sha512`
    * _Default value:_ `sha1`, `sha256`, `sha384`, `sha512`
* `tunnel2_phase1_lifetime_seconds` - (Optional, Integer) The lifetime for phase 1 of the IKE negotiation for the second VPN tunnel, in seconds.
    * _Valid values:_ From `900` to `28800`
    * _Default value:_ `28800`
* `tunnel2_phase2_dh_group_numbers` - (Optional, Set of integers) The Diffie-Hellman group numbers permitted for the second VPN tunnel for phase 2 of the IKE negotiation.
    * _Valid values:_ `0`, `2`, `5`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`
    * _Default value:_ `14`
    * _Constraints:_ No more than one value, `19`, `20` and `21` are not supported by `ikev1`
* `tunnel2_phase2_encryption_algorithms` - (Optional, Set of strings) The encryption algorithms permitted for the second VPN tunnel for phase 2 of the IKE negotiation.
    * _Valid values:_ `aes128`, `aes256`, `aes_ccm128`, `aes_ccm256`, `aes_ctr128`, `aes_ctr256`, `aes_gcm128`, `aes_gcm256`, `camellia128`, `camellia256`, `chacha20poly1305`
    * _Default value:_ All algorithms supported by the selected IKE version
    * _Constraints:_ `chacha20poly1305` is not supported by `ikev1`
* `tunnel2_phase2_integrity_algorithms` - (Optional, Set of strings) The integrity algorithms permitted for the second VPN tunnel for phase 2 of the IKE negotiation.
    * _Valid values:_ `sha1`, `sha256`, `sha384`, `sha512`
    * _Default value:_ `sha1`, `sha256`, `sha384`, `sha512`
* `tunnel2_phase2_lifetime_seconds` - (Optional, Integer) The lifetime for phase 2 of the IKE negotiation for the second VPN tunnel, in seconds.
    * _Valid values:_ From `900` to `3600`
    * _Default value:_ `3600`
    * _Constraints:_ Cannot exceed the value of `tunnel2_phase1_lifetime_seconds`
* `tunnel2_replay_window_size` - (Optional, Integer) The number of packets in an IKE replay window for the second VPN tunnel.
    * _Valid values:_ From `32` to `2048`
    * _Default value:_ `1024`

### Unsupported arguments

~> **Note** The platform does not support these arguments, so a request with any of them is rejected.
They may be present in the `terraform.tfstate` file with preset values.

The following arguments are not currently supported:

`enable_acceleration`, `local_ipv6_network_cidr`, `remote_ipv6_network_cidr`, `transit_gateway_id`, `tunnel1_dpd_timeout_action`, `tunnel2_dpd_timeout_action`, `tunnel1_dpd_timeout_seconds`, `tunnel2_dpd_timeout_seconds`, `tunnel1_inside_ipv6_cidr`, `tunnel2_inside_ipv6_cidr`, `tunnel1_rekey_fuzz_percentage`, `tunnel2_rekey_fuzz_percentage`, `tunnel1_rekey_margin_time_seconds`, `tunnel2_rekey_margin_time_seconds`, `tunnel1_startup_action`, `tunnel2_startup_action`.

The following arguments accept only one value, which is also their default:

* `static_routes_only` - `false`.
  Routing through static routes only is not supported.
* `tunnel_inside_ip_version` - `ipv4`.
  IPv6 traffic in the tunnels is not supported.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The ARN of the VPN connection.
* `customer_gateway_configuration` - (String, Sensitive) The configuration of the customer gateway of the VPN connection in the native XML format.
* `id` - (String) The ID of the VPN connection.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.
* `tunnel1_address` - (String) The public IP address of the first VPN tunnel.
* `tunnel1_bgp_asn` - (String) The BGP ASN of the first VPN tunnel (the VPN gateway side).
* `tunnel1_bgp_holdtime` - (Integer) The BGP hold time of the first VPN tunnel, in seconds.
* `tunnel1_cgw_inside_address` - (String) The link-local address of the first VPN tunnel (the customer gateway side).
* `tunnel1_vgw_inside_address` - (String) The link-local address of the first VPN tunnel (the VPN gateway side).
* `tunnel1_preshared_key` - (String, Sensitive) The preshared key of the first VPN tunnel.
* `tunnel2_address` - (String) The public IP address of the second VPN tunnel.
* `tunnel2_bgp_asn` - (String) The BGP ASN of the second VPN tunnel (the VPN gateway side).
* `tunnel2_bgp_holdtime` - (Integer) The BGP hold time of the second VPN tunnel, in seconds.
* `tunnel2_cgw_inside_address` - (String) The link-local address of the second VPN tunnel (the customer gateway side).
* `tunnel2_vgw_inside_address` - (String) The link-local address of the second VPN tunnel (the VPN gateway side).
* `tunnel2_preshared_key` - (String, Sensitive) The preshared key of the second VPN tunnel.
* `vgw_telemetry` - ([Block](#vgw_telemetry)) The telemetry of the VPN tunnels.

### vgw_telemetry

Each block has the following structure:

* `accepted_route_count` - (Integer) The number of accepted routes.
* `last_status_change` - (String) The date and time of the last change in the status.
* `outside_ip_address` - (String) The Internet-routable IP address of the outside interface of the VPN gateway.
* `status` - (String) The status of the VPN tunnel.
    * _Valid values:_ `UP`, `DOWN`
* `status_message` - (String) The description of the tunnel state.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`core_network_arn`, `core_network_attachment_arn`, `routes`, `transit_gateway_attachment_id`, `vgw_telemetry.certificate_arn`.

## Timeouts

Timeouts usage for VPN connections is not currently supported.

## Import

In Terraform v1.5.0 or later, the VPN connection can be imported by `id` using the `import` block.

```terraform
import {
  to = aws_vpn_connection.example
  id = "vpn-12345678"
}
```

In older Terraform versions, the VPN connection can be imported by its `id` using `terraform import`, for example:

```console
% terraform import aws_vpn_connection.example vpn-12345678
```
