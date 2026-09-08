---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_instances"
description: |-
  Provides lists of instance IDs, private IPs, and public IPs.
---

[describe-instances]: https://docs.k2.cloud/en/api/ec2/actions/instances/DescribeInstances.html
[outputs]: https://developer.hashicorp.com/terraform/language/values/outputs
[remote state]: https://developer.hashicorp.com/terraform/language/state/remote
[terraform_remote_state]: https://developer.hashicorp.com/terraform/language/state/remote-state-data

# Data Source: aws_instances

Provides lists of instance IDs, private IPs, and public IPs.

-> **Note:** It's a best practice to expose instance details via [outputs] or [remote state].
**Use the [`terraform_remote_state`][terraform_remote_state] data source instead** if you manage referenced instances via Terraform.

## Example usage

```terraform
data "aws_instances" "selected" {
  instance_tags = {
    type = "test"
  }

  filter {
    name   = "instance.group-id"
    values = ["sg-12345678"]
  }

  instance_state_names = ["running", "stopped"]
}

resource "aws_eip" "example" {
  count    = length(data.aws_instances.selected.ids)
  instance = data.aws_instances.selected.ids[count.index]
}
```

## Argument reference

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in the [EC2 API documentation][describe-instances]
* `instance_state_names` - (Optional, Set of strings) List of instance states that should be applicable to the desired instances.
    * _Valid values:_ `pending`, `running`, `shutting-down`, `stopped`, `stopping`, `terminated`
* `instance_tags` - (Optional, Map of strings) Key-value pairs. Must exactly match pairs on the desired resources.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The region.
* `ids` - (List of strings) The IDs of the instances found through the filter.
* `private_ips` - (List of strings) The private IP addresses of the instances found through the filter.
* `public_ips` - (List of strings) The public IP addresses of the instances found through the filter.
