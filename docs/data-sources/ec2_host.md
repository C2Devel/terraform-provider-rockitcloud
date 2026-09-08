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

## Example usage

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

## Argument reference

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in the [EC2 API documentation][describe-hosts]
* `host_id` - (Optional, String) The ID of the dedicated host.
* `tags` - (Optional, Map of strings) Key-value pairs. Must exactly match pairs on the desired resource.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the dedicated host.
* `auto_placement` - (String) Indicates whether automated placement is on or off.
* `availability_zone` - (String) The availability zone of the dedicated host.
* `cores` - (Integer) The number of cores on the dedicated host.
* `host_recovery` - (String) Indicates whether host recovery is enabled or disabled for the dedicated host.
* `id` - (String) The ID of the dedicated host.
* `instance_family` - (List of strings) The instance family supported by the dedicated host.
* `instance_type` - (String) The instance type supported by the dedicated host.
* `owner_id` - (String) The ID of the project the dedicated host belongs to.
* `sockets` - (Integer) The number of sockets on the dedicated host.
* `tags` - (Map of strings) Key-value pairs assigned to the resource.
* `total_vcpus` - (Integer) The total number of vCPUs on the dedicated host.
