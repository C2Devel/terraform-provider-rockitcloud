---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_ec2_host"
description: |-
  Provides information about a dedicated host.
---

[describe-hosts]: https://docs.k2.cloud/en/api/ec2/actions/hosts/DescribeHosts.html

# Data Source: aws_ec2_host

Provides information about a dedicated host.

## Example Usage

```terraform
data "aws_ec2_host" "selected" {
  host_id = aws_ec2_host.test.id
}
```

### Filter

```terraform
data "aws_ec2_host" "selected" {
  filter {
    name   = "auto-placement"
    values = ["on"]
  }

  filter {
    name   = "state"
    values = ["available"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-hosts]
* `host_id` - (Optional) The ID of the dedicated host.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the dedicated host.
* `auto_placement` - Indicates whether automated placement is on or off.
* `availability_zone` - Availability zone of the dedicated host.
* `cores` - Number of cores on the dedicated host.
* `host_recovery` - Indicates whether host recovery is enabled or disabled for the dedicated host.
* `id` - The ID of the dedicated host.
* `instance_family` - Instance family supported by the dedicated host.
* `owner_id` - The ID of the project that owns the dedicated host.
* `sockets` - Number of sockets on the dedicated host.
* `total_vcpus` - Total number of vCPUs on the dedicated host.
