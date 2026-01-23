---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_interface"
description: |-
  Provides information about a network interface.
---

[describe-network-interfaces]: https://docs.k2.cloud/en/api/ec2/actions/network_interfaces/DescribeNetworkInterfaces.html

# aws_network_interface

Provides information about a network interface.

## Example Usage

```terraform
data "aws_network_interface" "example" {
  id = "eni-xxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-network-interfaces]
* `id` - (Optional) The ID of the network interface.

## Attribute Reference

### Supported attributes

See the [`aws_network_interface`](network_interface.md) for details on the returned attributes.

Additionally, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the network interface.
* `association` - The association information for an Elastic IP address (IPv4) associated with the network interface.
  The structure of this block is [described below](#association).
* `availability_zone` - The availability zone.
* `description` - Description of the network interface.
* `mac_address` - The MAC address.
* `owner_id` - The project ID.
* `private_dns_name` - The private DNS name.
* `private_ip` - The private IPv4 address of the network interface within the subnet.
* `private_ips` - The private IPv4 addresses associated with the network interface.
* `security_groups` - The list of security groups for the network interface.
* `subnet_id` - The ID of the subnet.
* `tags` - Map of tags assigned to the network interface.
* `vpc_id` - The ID of the VPC.

#### association

* `allocation_id` - The allocation ID.
* `association_id` - The association ID.
* `customer_owned_ip` - The customer-owned IP address.
* `ip_owner_id` - The ID of the elastic IP address owner.
* `public_dns_name` - The public DNS name.
* `public_ip` - The address of the elastic IP address bound to the network interface.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`association.carrier_ip`, `interface_type`, `ipv6_addresses`, `outpost_arn`, `requester_id`.
