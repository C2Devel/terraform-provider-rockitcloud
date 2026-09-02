---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "aws_lb"
description: |-
  Manages a load balancer.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[elb]: https://docs.k2.cloud/en/services/elb/overview.html
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts

# Resource: aws_lb

Manages a load balancer.
For details about load balancers, see the [user documentation][elb].

~> **Important** An internet-facing load balancer requires the VPC to have an attached internet gateway and a default route (`0.0.0.0/0`) to it.
It's recommended to specify the route as an explicit dependency via `depends_on`.

## Example Usage

### Internal Application Load Balancer

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_subnet" "example" {
  vpc_id     = aws_vpc.example.id
  cidr_block = "10.1.1.0/24"

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_lb" "alb" {
  name               = "tf-alb"
  internal           = true
  load_balancer_type = "application"
  subnets            = [aws_subnet.example.id]

  tags = {
    Name = "tf-alb"
  }
}
```

### Internet-Facing Network Load Balancer

~> **Note** This example uses the VPC and subnet defined in the [Internal Application Load Balancer example](#internal-application-load-balancer).

```terraform
resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id

  tags = {
    Name = "tf-igw"
  }
}

resource "aws_route" "default_route" {
  route_table_id         = aws_vpc.example.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.example.id
}

resource "aws_eip" "example" {
  tags = {
    Name = "tf-eip"
  }
}

resource "aws_lb" "nlb" {
  depends_on = [aws_route.default_route]

  name               = "tf-nlb"
  internal           = false
  load_balancer_type = "network"

  subnet_mapping {
    subnet_id     = aws_subnet.example.id
    allocation_id = aws_eip.example.id
  }

  tags = {
    Name = "tf-nlb"
  }
}
```

## Argument Reference

The following arguments are supported:

* `internal` - (Optional) Indicates whether the load balancer will be internal or internet-facing.
* `load_balancer_type` - (Optional) The type of the load balancer.
    * _Valid values:_ `application`, `network`
* `name` - (Optional) The name of the load balancer.
    * _Value length:_ From 1 to 32 symbols
    * _Constraints:_
        * `name` cannot be specified if `name_prefix` is set
        * The value can contain only Latin letters, numbers, and hyphens (`-`)
        * The value must start and end with a Latin letter or number
        * The value cannot start with the prefix `internal-`
* `name_prefix` - (Optional) Creates a unique name beginning with the specified prefix.
    * _Value length:_ From 1 to 6 symbols
    * _Constraints:_
        * `name_prefix` cannot be specified if `name` is set
        * The value constraints are the same as for `name`

-> **Note** If `name` and `name_prefix` are not specified, Terraform will autogenerate a name with the prefix `tf-lb`.

* `subnet_mapping` - (Optional, Editable) List of subnet-ID-to-IP-address mappings.
  The structure of this block is [described below](#subnet_mapping).
    * _Constraints:_ `subnet_mapping` is required if the `subnets` argument is not specified
* `subnets` - (Optional, Editable) List of subnet IDs.
    * _Constraints:_
        * The `subnets` argument is required if `subnet_mapping` is not specified
        * All subnets must be from different availability zones

~> **Note** You can only add new subnets to the `subnets` or `subnet_mapping` list, subnets cannot be removed.
  
* `tags` - (Optional, Editable) Map of tags to assign to the load balancer.
  If a provider [`default_tags` configuration block][default-tags] is used,
  tags with matching keys will overwrite those defined at the provider level.

### subnet_mapping

The `subnet_mapping` block has the following structure:

* `subnet_id` - (Required, Editable) The ID of the subnet.
* `allocation_id` - (Optional, Editable) The ID of the Elastic IP address allocation.
  The _internet-facing_ load balancer will be available at this IP address.
* `private_ipv4_address` - (Optional, Editable) The private IP address within the specified subnet.
  The _internal_ load balancer will be available at this IP address.

~> **Note** All subnets specified in the `subnet_mapping` blocks must be from different availability zones.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the load balancer.
* `dns_name` - The DNS name of the load balancer.
* `id` - The ARN of the load balancer.
* `tags_all` - Map of tags assigned to the load balancer,
  including those inherited from the provider [`default_tags` configuration block][default-tags].
* `vpc_id` - The ID of the VPC.
* `zone_id` - The ID of the Route53 hosted zone associated with the load balancer.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`access_logs`, `arn_suffix`, `customer_owned_ipv4_pool`, `desync_mitigation_mode`, `drop_invalid_header_fields`, `enable_cross_zone_load_balancing`, `enable_deletion_protection`, `enable_http2`, `enable_waf_fail_open`, `idle_timeout`, `ip_address_type`, `security_groups`, `subnet_mapping.ipv6_address`, `subnet_mapping.outpost_id`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

- `create` - (Default `10 minutes`) Used when creating the load balancer.
- `update` - (Default `10 minutes`) Used when updating the load balancer.
- `delete` - (Default `10 minutes`) Used when destroying the load balancer.

## Import

The load balancer can be imported using `arn`, e.g.,

```
$ terraform import aws_lb.alb arn:c2:elasticloadbalancing::project-name@customer-name:loadbalancer/app/lb-12345678
```
