---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_subnet_ids"
description: |-
  Provides a list of subnet IDs for a VPC.
---

[describe-subnets]: https://docs.k2.cloud/en/api/ec2/actions/subnets/DescribeSubnets.html

# Data Source: aws_subnet_ids

Provides a list of subnet IDs for a VPC.

## Example Usage

The following shows all CIDR blocks for every subnet ID in a VPC.

```terraform
variable vpc_id {}

data "aws_subnets" "example" {
  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }
}

data "aws_subnet" "example" {
  for_each = toset(data.aws_subnets.example.ids)
  id       = each.value
}

output "subnet_cidr_blocks" {
  value = [for s in data.aws_subnet.example : s.cidr_block]
}
```

The following example retrieves a set of all subnets in a VPC with a custom
tag of `Tier` set to a value of "Private" so that the `aws_instance` resource
can loop through the subnets, putting instances across availability zones.

```terraform
variable vpc_id {}

data "aws_subnets" "private" {
  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }

  tags = {
    Tier = "Private"
  }
}

resource "aws_instance" "app" {
  for_each      = toset(data.aws_subnets.example.ids)
  ami           = "cmi-12345678" # add image id, change instance type if needed
  instance_type = "m1.micro"
  subnet_id     = each.value
}
```

For matching against tag `Name`, use:

```terraform
data "aws_subnet_ids" "selected" {
  filter {
    name   = "tag:Name"
    values = [""] # insert values here
  }
}
```

## Argument Reference

* `vpc_id` - (Required) The VPC ID that you want to filter from.
* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-subnets]
* `tags` - (Optional) Map of tags, each pair of which must exactly match
  a pair on the desired subnets.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `ids` - Set of found subnet IDs.
