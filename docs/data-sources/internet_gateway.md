---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_internet_gateway"
description: |-
  Provides information about an internet gateway.
---

[describe-igws]: https://docs.k2.cloud/en/api/ec2/actions/internet_gateways/DescribeInternetGateways.html

# Data Source: aws_internet_gateway

Provides information about an internet gateway.

## Example usage

### Basic example

```terraform
data "aws_internet_gateway" "selected" {
  filter {
    name   = "attachment.vpc-id"
    values = ["vpc-12345678"]
  }
}
```

## Argument reference

The arguments of this data source act as filters for querying the available internet gateway.

~> **Note** The given filters must exactly match the resource whose data will be exported as attributes.

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-igws]
* `internet_gateway_id` - (Optional, String) The ID of the internet gateway.
* `tags` - (Optional, Map of strings) Key-value pairs.
    Must exactly match pairs on the required resource.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

If any fields are missing from the configuration, then this data source will populate them with data for the selected internet gateway.

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the internet gateway.
* `attachments` - ([Block](#attachments)) The list of VPC attachments to the internet gateway. The list can contain 0 or 1 element.
* `id` - (String) The ID of the internet gateway.
* `owner_id` - (String) The ID of the project that owns the internet gateway.

### attachments

* `state` - (String) The current state of the VPC attachment to the internet gateway.
    * _Valid values:_ `attached`, `attaching`, `available`, `detached`, `detaching`
* `vpc_id` - (String) The ID of the attached VPC.
