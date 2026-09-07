---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_subnet"
description: |-
  Manages a subnet.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[subnets]: https://docs.k2.cloud/en/services/networking/subnets.html
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_subnet

Manages a subnet.

For more information about subnets, see the documentation on [Subnets][subnets].

## Example usage

### Basic example

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_subnet" "example" {
  vpc_id     = aws_vpc.example.id
  cidr_block = "10.1.1.0/24"

  tags = {
    Name = "Main"
  }
}
```

## Argument reference

The following arguments are supported:

* `vpc_id` - (Required, Forces new resource, String) The ID of the VPC.
* `availability_zone` - (Optional, Forces new resource, String) The availability zone for the subnet.
* `cidr_block` - (Required, Forces new resource, String) The IPv4 CIDR block for the subnet.
* `map_public_ip_on_launch` - (Optional, Editable, Boolean, **Deprecated**) Indicates whether public IP addresses will be associated with instances created in this subnet. Addresses are associated only if there are available allocated Elastic IP addresses.
    * _Default value:_ `false`
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource. If the [`default_tags` configuration block][default-tags] block is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

~> **Note** The argument `map_public_ip_on_launch` is deprecated and will be removed in a future release.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The ID of the subnet.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`arn`, `assign_ipv6_address_on_creation`, `availability_zone_id`, `customer_owned_ipv4_pool`, `enable_dns64`, `enable_resource_name_dns_a_record_on_launch`, `enable_resource_name_dns_aaaa_record_on_launch`, `ipv6_cidr_block`, `ipv6_cidr_block_association_id`, `ipv6_native`, `map_customer_owned_ip_on_launch`, `outpost_arn`, `owner_id`, `private_dns_hostname_type_on_launch`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `10m`) How long to wait for a subnet to be created.
- `delete` - (Default `20m`) How long to wait for a subnet to be deleted.

## Import

Subnets can be imported using the `id`, for example:

```
$ terraform import aws_subnet.public_subnet subnet-12345678
```
