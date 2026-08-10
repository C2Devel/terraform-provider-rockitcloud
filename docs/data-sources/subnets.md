---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_subnets"
description: |-
  Provides a list of subnet IDs.
---

[describe-subnets]: https://docs.k2.cloud/en/api/ec2/actions/subnets/DescribeSubnets.html

# Data Source: aws_subnets

Provides a list of subnet IDs.

## Example usage

### Specific example: get all CIDR blocks for every subnet ID in a VPC

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

### Specific example: retrieve a set of all subnets in a VPC

The following example retrieves a set of all subnets in a VPC with a custom `Tier` tag set to `Private`.
So the `aws_instance` resource can loop through the subnets, putting instances across availability zones.

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
data "aws_subnets" "selected" {
  filter {
    name   = "tag:Name"
    values = [""] # insert values here
  }
}
```

## Argument reference

The following arguments are supported:

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-subnets]
* `tags` - (Optional, Map of strings) Key-value pairs.
    Must exactly match pairs on the required resources.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attribute is exported:

* `id` - (String) The region.
* `ids` - (List of strings) The list of all the subnet IDs found.
