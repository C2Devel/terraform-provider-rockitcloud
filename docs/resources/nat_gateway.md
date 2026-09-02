---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_nat_gateway"
description: |-
  Manages a regional NAT gateway.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block

# Resource: aws_nat_gateway

Manages a regional NAT gateway.
A regional NAT gateway is a single NAT gateway that works across all availability zones in your VPC.

-> **Note** Only one NAT gateway in the `pending` or `available` state is allowed for each VPC.

~> **Important** The VPC must have an attached internet gateway.
It's recommended to specify it as an explicit dependency via `depends_on`.

## Example usage

### Basic example

```terraform
resource "aws_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_nat_gateway" "example" {
  vpc_id = aws_vpc.vpc.id

  tags = {
    Name = "tf-nat-gw"
  }

  depends_on = [aws_internet_gateway.igw]
}

resource "aws_route" "default_route" {
  route_table_id         = aws_vpc.vpc.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = aws_nat_gateway.example.id
}
```

### Specific example: manually specified Elastic IP addresses

```terraform
resource "aws_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.vpc.id
}

resource "aws_eip" "eip" {}

resource "aws_nat_gateway" "example" {
  vpc_id = aws_vpc.vpc.id

  availability_zone_addresses {
    availability_zone = "ru-msk-vol51"
    allocation_id     = aws_eip.eip.id
  }

  depends_on = [aws_internet_gateway.igw]
}

resource "aws_route" "default_route" {
  route_table_id         = aws_vpc.vpc.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = aws_nat_gateway.example.id
}
```

## Argument reference

The following arguments are required:

* `vpc_id` - (Required, Forces new resource, String) The ID of the VPC in which the NAT gateway is created.

The following arguments are optional:

* `availability_mode` - (Optional, Forces new resource, String) The availability mode of the NAT gateway.
    * _Valid values:_ `regional`
    * _Default value:_ `regional`
* `availability_zone_addresses` - (Optional, Editable, [Block](#availability_zone_addresses)) The Elastic IP addresses to use for handling outbound NAT traffic in specific availability zones.
  If the block isn't specified, the Elastic IP addresses are allocated and managed automatically.
    * _Constraints:_ Automatically allocated addresses cannot be managed via `availability_zone_addresses`, so adding the block to such a NAT gateway forces a new resource
* `connectivity_type` - (Optional, Forces new resource, String) The connectivity type of the NAT gateway.
    * _Valid values:_ `public`
    * _Default value:_ `public`
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource.
  If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

### availability_zone_addresses

The block has the following structure:

* `allocation_id` - (Required, Editable, String) The allocation ID of the Elastic IP address to use in the availability zone.
* `availability_zone` - (Required, Editable, String) The name of the availability zone.

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `auto_provision_zones` - (String) The state of automatic Elastic IP address allocation for each availability zone.
    * _Valid values:_ `disabled`, `enabled`
* `id` - (String) The ID of the NAT gateway.
* `nat_gateway_addresses` - ([Block](#nat_gateway_addresses)) The set of addresses of the NAT gateway.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### nat_gateway_addresses

Each block has the following structure:

* `allocation_id` - (String) The allocation ID of the Elastic IP address.
* `association_id` - (String) The association ID of the Elastic IP address.
* `availability_zone` - (String) The availability zone of the Elastic IP address.
* `is_primary` - (Boolean) Indicates whether the Elastic IP address is the primary address of the NAT gateway.
* `private_ip` - (String) The private IP address.
* `public_ip` - (String) The public IP address.
* `status` - (String) The status of the Elastic IP address.

## Timeouts

Timeouts usage for NAT gateways is not currently supported.

## Import

In Terraform v1.5.0 or later, the NAT gateway can be imported by `id` using the `import` block.

```terraform
import {
  to = aws_nat_gateway.example
  id = "nat-12345678"
}
```

In older Terraform versions, the NAT gateway can be imported by its `id` using `terraform import`, for example:

```console
% terraform import aws_nat_gateway.example nat-12345678
```
