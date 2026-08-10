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

## Example usage

```terraform
resource "aws_placement_group" "example" {
  name     = "test-pg"
  strategy = "cluster"
}
```

## Argument reference

The following arguments are supported:

* `name` - (Required, Forces new resource, String) The name of the placement group.
* `strategy` - (Required, Forces new resource, String) The placement strategy.
    * _Valid values:_ `cluster`, `partition`, `spread`
* `partition_count` - (Optional, Forces new resource, Integer) The number of partitions to create in the placement group.
    * _Constraints:_ Must be between 0 and 7. Can only be set when `strategy` is set to `partition`.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the placement group. If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the placement group.
* `id` - (String) The name of the placement group.
* `placement_group_id` - (String) The ID of the placement group.
* `tags_all` - (Map of strings) Key-value pairs assigned to the placement group, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

## Timeouts

Timeouts usage for the placement group is not currently supported.

## Import

Placement groups can be imported using `id`, for example:

```
$ terraform import aws_placement_group.prod_pg production-placement-group
```