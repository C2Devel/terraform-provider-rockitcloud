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
[tunnel-options]: https://docs.k2.cloud/en/api/ec2/datatypes/VpnTunnelOptionsSpecification.html

# Resource: aws_vpn_connection

Manages a Site-to-Site VPN connection.
A Site-to-Site VPN connection is an IPsec VPN connection between a VPC and an on-premises network.

For more information, see the documentation on [VPN connections][vpn-connections].
For more information about the tunnel options and the compatibility of the algorithms with the IKE versions, see [EC2 API documentation][tunnel-options].

~> **Note** All arguments including `tunnel1_preshared_key` and `tunnel2_preshared_key` are stored in the raw state as plain text.
[Read more about sensitive data in state][sensitive-data].

-> **Note** A VPN gateway is created automatically for each VPC.
For convenience, the ID of the VPN gateway is the same as the ID of the VPC to which it belongs (`vpc-ABCD1234`/`vgw-ABCD1234`).
Either of them can be specified in the `vpn_gateway_id` argument.

~> **Important** Only one VPN connection can exist between a customer gateway and a VPN gateway.

~> **Important** The number of VPN tunnels depends on the specified tunnel options.
If no tunnel options are specified or the options of both tunnels are specified, the VPN connection has two tunnels in different availability zones.
If only the options of the first tunnel are specified, the VPN connection has a single tunnel.
The options of the second tunnel can be specified only together with at least one option of the first tunnel.
Two tunnels require a region with at least two availability zones.
Removing all tunnel options from the configuration does not change the number of tunnels of the existing VPN connection.

~> **Note** A VPN connection cannot be modified, except for its tags.
A change of any other argument forces a new resource, and the addresses of the tunnels change.

## Example usage

### Basic example

By default, a high-availability VPN connection is created.
It has two tunnels, which are terminated in different availability zones.
Such a VPN connection supports dynamic routing over BGP only.
To install the routes received from the customer gateway into a route table, use the [`aws_vpn_gateway_route_propagation`](vpn_gateway_route_propagation.md) resource.

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

The options of each tunnel are specified separately.

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

  tunnel1_ike_versions                 = ["ikev2"]
  tunnel1_phase1_encryption_algorithms = ["aes256", "aes_gcm256"]
  tunnel1_phase1_integrity_algorithms  = ["sha256", "sha512"]
  tunnel1_phase1_dh_group_numbers      = [19, 20, 21]
  tunnel1_phase2_dh_group_numbers      = [19]
  tunnel1_replay_window_size           = 512

  tunnel2_ike_versions                 = ["ikev2"]
  tunnel2_phase1_encryption_algorithms = ["aes256", "aes_gcm256"]
  tunnel2_phase1_integrity_algorithms  = ["sha256", "sha512"]
  tunnel2_phase1_dh_group_numbers      = [19, 20, 21]
  tunnel2_phase2_dh_group_numbers      = [19]
  tunnel2_replay_window_size           = 512
}
```

### Specific example: VPN connection with a single tunnel and static routing

A VPN connection has a single tunnel if only the options of the first tunnel are specified.
Only such a VPN connection can be a target of static routes.

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

  tunnel1_ike_versions = ["ikev2"]
}

resource "aws_route" "example" {
  route_table_id         = aws_vpc.example.main_route_table_id
  destination_cidr_block = "192.168.0.0/24"
  gateway_id             = aws_vpn_connection.example.id
}
```

## Argument reference

The following arguments are required:

* `customer_gateway_id` - (Required, Forces new resource, String) The ID of the customer gateway.
* `type` - (Required, Forces new resource, String) The type of the VPN connection.
    * _Valid values:_ `ipsec.1` or `ipsec.legacy`
    * _Constraints:_ Must match the type of the customer gateway
* `vpn_gateway_id` - (Required, Forces new resource, String) The ID of the VPN gateway or of the VPC to which it belongs.
  The exported value is always in the `vgw-ABCD1234` form.

The following arguments are optional:

* `local_ipv4_network_cidr` - (Optional, Forces new resource, String) The IPv4 CIDR on the customer gateway (on-premises) side of the VPN connection.
    * _Default value:_ `0.0.0.0/0`
    * _Constraints:_ Must not fall within the `169.254.0.0/16` range
* `remote_ipv4_network_cidr` - (Optional, Forces new resource, String) The IPv4 CIDR on the cloud side of the VPN connection.
    * _Default value:_ `0.0.0.0/0`
    * _Constraints:_ Must not fall within the `169.254.0.0/16` range
* `static_routes_only` - (Optional, Forces new resource, Boolean) Indicates whether the VPN connection uses static routes only.
  Currently, the argument is not used.
    * _Default value:_ `false`
    * _Valid values:_ `false`
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the VPN connection.
  If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
* `tunnel1_ike_versions` - (Optional, Forces new resource, Set of strings) The IKE version permitted for the first VPN tunnel.
    * _Default value:_ `ikev1`
    * _Valid values:_ `ikev1` or `ikev2`
    * _Constraints:_ Only one IKE version can be specified
* `tunnel1_inside_cidr` - (Optional, Forces new resource, String) The CIDR block of the inside IP addresses for the first VPN tunnel.
  The first usable address of the block is assigned to the VPN gateway and the second one to the customer gateway.
  To assign another address to the VPN gateway, specify it instead of the network address.
    * _Default value:_ The block is assigned automatically
    * _Constraints:_
        * A `/30` CIDR block from the `169.254.252.0/22` range
        * The block must be unique among the VPN connections of the VPN gateway
    * _Example:_ `169.254.252.2/30`
* `tunnel1_phase1_dh_group_numbers` - (Optional, Forces new resource, Set of integers) The Diffie-Hellman group numbers permitted for the first VPN tunnel for phase 1 of the IKE negotiation.
    * _Default value:_ `5`, `14`, `15`, `16`, `17`, `18`
    * _Valid values:_ `5`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`
    * _Constraints:_ No more than six values
* `tunnel1_phase1_encryption_algorithms` - (Optional, Forces new resource, Set of strings) The encryption algorithms permitted for the first VPN tunnel for phase 1 of the IKE negotiation.
    * _Default value:_ All algorithms supported by the selected IKE version
    * _Valid values:_ `aes128`, `aes256`, `aes_ctr128`, `aes_ctr256`, `aes_gcm128`, `aes_gcm256`, `camellia128`, `camellia256`, `chacha20poly1305`
    * _Constraints:_ `aes_gcm128`, `aes_gcm256` and `chacha20poly1305` are not supported by `ikev1`
* `tunnel1_phase1_integrity_algorithms` - (Optional, Forces new resource, Set of strings) The integrity algorithms permitted for the first VPN tunnel for phase 1 of the IKE negotiation.
    * _Default value:_ `sha1`, `sha256`, `sha384`, `sha512`
    * _Valid values:_ `sha1`, `sha256`, `sha384`, `sha512`
* `tunnel1_phase1_lifetime_seconds` - (Optional, Forces new resource, Integer) The lifetime for phase 1 of the IKE negotiation for the first VPN tunnel, in seconds.
    * _Default value:_ `28800`
    * _Valid values:_ From `900` to `28800`
* `tunnel1_phase2_dh_group_numbers` - (Optional, Forces new resource, Set of integers) The Diffie-Hellman group number permitted for the first VPN tunnel for phase 2 of the IKE negotiation.
  The `0` value disables Perfect Forward Secrecy (PFS), which is not recommended.
    * _Default value:_ `14`
    * _Valid values:_ `0`, `5`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`
    * _Constraints:_
        * Only one value can be specified
        * `19`, `20` and `21` are not supported by `ikev1`
* `tunnel1_phase2_encryption_algorithms` - (Optional, Forces new resource, Set of strings) The encryption algorithms permitted for the first VPN tunnel for phase 2 of the IKE negotiation.
    * _Default value:_ All algorithms supported by the selected IKE version
    * _Valid values:_ `aes128`, `aes256`, `aes_ccm128`, `aes_ccm256`, `aes_ctr128`, `aes_ctr256`, `aes_gcm128`, `aes_gcm256`, `camellia128`, `camellia256`, `chacha20poly1305`
    * _Constraints:_ `chacha20poly1305` is not supported by `ikev1`
* `tunnel1_phase2_integrity_algorithms` - (Optional, Forces new resource, Set of strings) The integrity algorithms permitted for the first VPN tunnel for phase 2 of the IKE negotiation.
    * _Default value:_ `sha1`, `sha256`, `sha384`, `sha512`
    * _Valid values:_ `sha1`, `sha256`, `sha384`, `sha512`
* `tunnel1_phase2_lifetime_seconds` - (Optional, Forces new resource, Integer) The lifetime for phase 2 of the IKE negotiation for the first VPN tunnel, in seconds.
    * _Default value:_ `3600`
    * _Valid values:_ From `900` to `3600`
    * _Constraints:_ Cannot exceed the value of `tunnel1_phase1_lifetime_seconds`
* `tunnel1_preshared_key` - (Optional, Forces new resource, String, Sensitive) The preshared key (PSK) of the first VPN tunnel.
    * _Default value:_ The key is generated automatically
    * _Constraints:_
        * From 8 to 64 characters
        * The value can contain only alphanumeric characters, periods (`.`) and underscores (`_`)
        * The value cannot start with zero (`0`)
* `tunnel1_replay_window_size` - (Optional, Forces new resource, Integer) The number of packets in an IKE replay window for the first VPN tunnel.
    * _Default value:_ `1024`
    * _Valid values:_ From `32` to `2048`
* `tunnel2_ike_versions` - (Optional, Forces new resource, Set of strings) The IKE version permitted for the second VPN tunnel.
    * _Default value:_ `ikev1`
    * _Valid values:_ `ikev1` or `ikev2`
    * _Constraints:_ Only one IKE version can be specified
* `tunnel2_inside_cidr` - (Optional, Forces new resource, String) The CIDR block of the inside IP addresses for the second VPN tunnel.
  The first usable address of the block is assigned to the VPN gateway and the second one to the customer gateway.
  To assign another address to the VPN gateway, specify it instead of the network address.
    * _Default value:_ The block is assigned automatically
    * _Constraints:_
        * A `/30` CIDR block from the `169.254.252.0/22` range
        * The block must be unique among the VPN connections of the VPN gateway
    * _Example:_ `169.254.252.2/30`
* `tunnel2_phase1_dh_group_numbers` - (Optional, Forces new resource, Set of integers) The Diffie-Hellman group numbers permitted for the second VPN tunnel for phase 1 of the IKE negotiation.
    * _Default value:_ `5`, `14`, `15`, `16`, `17`, `18`
    * _Valid values:_ `5`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`
    * _Constraints:_ No more than six values
* `tunnel2_phase1_encryption_algorithms` - (Optional, Forces new resource, Set of strings) The encryption algorithms permitted for the second VPN tunnel for phase 1 of the IKE negotiation.
    * _Default value:_ All algorithms supported by the selected IKE version
    * _Valid values:_ `aes128`, `aes256`, `aes_ctr128`, `aes_ctr256`, `aes_gcm128`, `aes_gcm256`, `camellia128`, `camellia256`, `chacha20poly1305`
    * _Constraints:_ `aes_gcm128`, `aes_gcm256` and `chacha20poly1305` are not supported by `ikev1`
* `tunnel2_phase1_integrity_algorithms` - (Optional, Forces new resource, Set of strings) The integrity algorithms permitted for the second VPN tunnel for phase 1 of the IKE negotiation.
    * _Default value:_ `sha1`, `sha256`, `sha384`, `sha512`
    * _Valid values:_ `sha1`, `sha256`, `sha384`, `sha512`
* `tunnel2_phase1_lifetime_seconds` - (Optional, Forces new resource, Integer) The lifetime for phase 1 of the IKE negotiation for the second VPN tunnel, in seconds.
    * _Default value:_ `28800`
    * _Valid values:_ From `900` to `28800`
* `tunnel2_phase2_dh_group_numbers` - (Optional, Forces new resource, Set of integers) The Diffie-Hellman group number permitted for the second VPN tunnel for phase 2 of the IKE negotiation.
  The `0` value disables Perfect Forward Secrecy (PFS), which is not recommended.
    * _Default value:_ `14`
    * _Valid values:_ `0`, `5`, `14`, `15`, `16`, `17`, `18`, `19`, `20`, `21`
    * _Constraints:_
        * Only one value can be specified
        * `19`, `20` and `21` are not supported by `ikev1`
* `tunnel2_phase2_encryption_algorithms` - (Optional, Forces new resource, Set of strings) The encryption algorithms permitted for the second VPN tunnel for phase 2 of the IKE negotiation.
    * _Default value:_ All algorithms supported by the selected IKE version
    * _Valid values:_ `aes128`, `aes256`, `aes_ccm128`, `aes_ccm256`, `aes_ctr128`, `aes_ctr256`, `aes_gcm128`, `aes_gcm256`, `camellia128`, `camellia256`, `chacha20poly1305`
    * _Constraints:_ `chacha20poly1305` is not supported by `ikev1`
* `tunnel2_phase2_integrity_algorithms` - (Optional, Forces new resource, Set of strings) The integrity algorithms permitted for the second VPN tunnel for phase 2 of the IKE negotiation.
    * _Default value:_ `sha1`, `sha256`, `sha384`, `sha512`
    * _Valid values:_ `sha1`, `sha256`, `sha384`, `sha512`
* `tunnel2_phase2_lifetime_seconds` - (Optional, Forces new resource, Integer) The lifetime for phase 2 of the IKE negotiation for the second VPN tunnel, in seconds.
    * _Default value:_ `3600`
    * _Valid values:_ From `900` to `3600`
    * _Constraints:_ Cannot exceed the value of `tunnel2_phase1_lifetime_seconds`
* `tunnel2_preshared_key` - (Optional, Forces new resource, String, Sensitive) The preshared key (PSK) of the second VPN tunnel.
    * _Default value:_ The key is generated automatically
    * _Constraints:_
        * From 8 to 64 characters
        * The value can contain only alphanumeric characters, periods (`.`) and underscores (`_`)
        * The value cannot start with zero (`0`)
* `tunnel2_replay_window_size` - (Optional, Forces new resource, Integer) The number of packets in an IKE replay window for the second VPN tunnel.
    * _Default value:_ `1024`
    * _Valid values:_ From `32` to `2048`
* `tunnel_inside_ip_version` - (Optional, Forces new resource, String) The IP version of the traffic inside the VPN tunnels.
  Currently, the argument is not used.
    * _Default value:_ `ipv4`
    * _Valid values:_ `ipv4`

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The ARN of the VPN connection.
* `customer_gateway_configuration` - (String, Sensitive) The configuration of the customer gateway of the VPN connection in the native XML format.
* `id` - (String) The ID of the VPN connection.
* `tags_all` - (Map of strings) Key-value pairs assigned to the VPN connection, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.
* `tunnel1_address` - (String) The public IP address of the first VPN tunnel.
* `tunnel1_bgp_asn` - (String) The BGP ASN of the first VPN tunnel (the VPN gateway side).
* `tunnel1_bgp_holdtime` - (Integer) The BGP hold time of the first VPN tunnel, in seconds.
* `tunnel1_cgw_inside_address` - (String) The link-local address of the first VPN tunnel (the customer gateway side).
* `tunnel1_vgw_inside_address` - (String) The link-local address of the first VPN tunnel (the VPN gateway side).
* `tunnel2_address` - (String) The public IP address of the second VPN tunnel.
* `tunnel2_bgp_asn` - (String) The BGP ASN of the second VPN tunnel (the VPN gateway side).
* `tunnel2_bgp_holdtime` - (Integer) The BGP hold time of the second VPN tunnel, in seconds.
* `tunnel2_cgw_inside_address` - (String) The link-local address of the second VPN tunnel (the customer gateway side).
* `tunnel2_vgw_inside_address` - (String) The link-local address of the second VPN tunnel (the VPN gateway side).
* `vgw_telemetry` - ([Block](#vgw_telemetry)) The telemetry of the VPN tunnels.

-> **Note** The attributes of the second tunnel are empty if the VPN connection has a single tunnel.

### vgw_telemetry

Each block has the following structure:

* `accepted_route_count` - (Integer) The number of accepted routes.
* `last_status_change` - (String) The date and time of the last change in the status.
* `outside_ip_address` - (String) The Internet-routable IP address of the outside interface of the VPN gateway.
* `status` - (String) The status of the VPN tunnel.
    * _Valid values:_ `DOWN` or `UP`
* `status_message` - (String) The description of the tunnel state.

## Timeouts

Timeouts usage for VPN connections is not currently supported.

## Import

In Terraform v1.5.0 or later, the VPN connection can be imported using `id` in the `import` block, for example:

```terraform
import {
  to = aws_vpn_connection.example
  id = "vpn-12345678"
}
```

In older Terraform versions, the VPN connection can be imported using `id`, for example:

```
$ terraform import aws_vpn_connection.example vpn-12345678
```
