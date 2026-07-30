---
subcategory: "Auto Scaling"
layout: "aws"
page_title: "aws_autoscaling_policy"
description: |-
  Manages an Auto Scaling policy.
---

# Resource: aws_autoscaling_policy

Manages an Auto Scaling policy.

~> **Note** You may want to omit the `desired_capacity` attribute from the attached `aws_autoscaling_group` when using Auto Scaling policies.
It's good practice to pick either manual or dynamic (policy-based) scaling.

## Example usage

```terraform
resource "aws_autoscaling_policy" "example" {
  name                   = "terraform-test"
  scaling_adjustment     = 4
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 300
  autoscaling_group_name = "example-asg" # asg is created manually
}
```

## Argument reference

The following arguments are required:

* `autoscaling_group_name` - (Required, Forces new resource, String) The name of the Auto Scaling group.
* `name` - (Required, Forces new resource, String) The name of the policy.

The following arguments are optional:

* `adjustment_type` - (Optional, Editable, String) Specifies whether the adjustment is an absolute number or a percentage of the current capacity.
    * _Valid values:_ `ChangeInCapacity`, `ExactCapacity`, `PercentChangeInCapacity`
* `cooldown` - (Optional, Editable, Integer) The amount of time in seconds, after a scaling activity completes and before the next scaling activity can start.
* `min_adjustment_magnitude` - (Optional, Editable, Integer) The minimum value to scale by when `adjustment_type` is set to `PercentChangeInCapacity`.
    * _Constraints:_ At least 1
* `policy_type` - (Optional, Editable, String) The policy type.
    * _Default value:_ `SimpleScaling`
    * _Valid values:_ `SimpleScaling`
* `scaling_adjustment` - (Optional, Editable, Integer) The amount by which the Auto Scaling group is scaled when the scaling policy is executed.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `adjustment_type` - (String) The adjustment type of the scaling policy.
* `arn` - (String) The Amazon Resource Name (ARN) of the scaling policy.
* `autoscaling_group_name` - (String) The name of the Auto Scaling group assigned to the scaling policy.
* `id` - (String) The name of the scaling policy.
* `name` - (String) The name of the scaling policy.
* `policy_type` - (String) The type of the scaling policy.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`estimated_instance_warmup`, `metric_aggregation_type`, `predictive_scaling_configuration`, `step_adjustment`, `target_tracking_configuration`.

## Timeouts

Timeouts usage for Auto Scaling policies is not currently supported.

## Import

The Auto Scaling policy can be imported using the `autoscaling_group_name` and `name` separated by `/`.

```
$ terraform import aws_autoscaling_policy.test-policy asg-name/policy-name
```
