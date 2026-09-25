---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_placement_group"
description: |-
  Manages an EC2 placement group.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[placement-groups]: https://docs.k2.cloud/en/services/compute/placementgroups.html

# Resource: aws_placement_group

Manages an EC2 placement group.
For more information, see the documentation on [placement groups][placement-groups].

## Example Usage

```terraform
resource "aws_placement_group" "example" {
  name     = "test-pg"
  strategy = "cluster"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the placement group.
* `strategy` - (Required) The placement strategy.
    * _Valid values:_ `spread`
* `tags` - (Optional) Map of tags to assign to the placement group. If a provider [`default_tags` configuration block][default-tags] is used, tags with matching keys will overwrite those defined at the provider level.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the placement group.
* `id` - The name of the placement group.
* `placement_group_id` - The ID of the placement group.
* `tags_all` - Map of tags assigned to the placement group, including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `partition_count`.

## Import

Placement groups can be imported using the `name`, e.g.,

```
$ terraform import aws_placement_group.prod_pg production-placement-group
```
