---
subcategory: "Auto Scaling"
layout: "aws"
page_title: "aws_autoscaling_group"
description: |-
  Provides information about an Auto Scaling group.
---

# Data Source: aws_autoscaling_group

Provides information about an Auto Scaling group.

## Example usage

```terraform
data "aws_autoscaling_group" "example" {
  name = "example-asg"
}
```

## Argument reference

* `name` - (Required, String) The name of the Auto Scaling group.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the Auto Scaling group.
* `availability_zones` - (Set of strings) One or more availability zones for the group.
* `default_cooldown` - (Integer) The amount of time in seconds after a scaling activity completes before another scaling activity can start.
* `desired_capacity` - (Integer) The desired size of the group.
* `health_check_grace_period` - (Integer) The amount of time in seconds, after which the Auto Scaling group can perform a health check on its instances.
* `id` - (String) The name of the Auto Scaling group.
* `max_size` - (Integer) The maximum size of the group.
* `min_size` - (Integer) The minimum size of the group.
* `name` - (String) The name of the Auto Scaling group.
* `new_instances_protected_from_scale_in` - (Boolean) Indicates whether new instances are protected from deletion when the Auto Scaling group is scaled in.
* `status` - (String) The status of the Auto Scaling group when it is deleted.
* `vpc_zone_identifier` - (String) The IDs of the subnets in which instances will be created.
* `launch_template` - (List of [Block](#launch_template)) The launch template for the group.

### launch_template

* `id` - (String) The ID of the launch template.
* `name` - (String) The name of the launch template.
* `version` - (String) The version of the launch template.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`enabled_metrics`, `health_check_type`, `launch_configuration`, `load_balancers`, `placement_group`, `service_linked_role_arn`, `target_group_arns`, `termination_policies`.
