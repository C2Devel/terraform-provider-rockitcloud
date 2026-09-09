---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "aws_lb_target_group"
description: |-
  Manages a target group.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[target-groups]: https://docs.k2.cloud/en/services/elb/target_groups.html

# Resource: aws_lb_target_group

Manages a target group.
For details about target groups, see the [user documentation][target-groups].

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_lb_target_group" "example" {
  name = "tf-lb-tg"

  target_type = "instance"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = aws_vpc.example.id

  tags = {
    Name = "tf-lb-tg"
  }
}
```

## Argument Reference

The following arguments are supported:

* `port` - (Required) The port on which targets will receive requests.
  Can be overridden when registering a specific target.
    * _Valid values:_ From 1 to 65535
* `protocol` - (Required) The protocol that will be used for routing traffic to the targets.
    * _Valid values:_ `TCP`, `UDP`, `HTTP`
* `vpc_id` - (Required) The ID of the VPC.
* `health_check` - (Optional, Editable) The health check configuration.
  The structure of this block is [described below](#health_check).
* `name` - (Optional) The name of the target group.
    * _Value length:_ From 1 to 32 symbols
    * _Constraints:_
        * `name` cannot be specified if `name_prefix` is set
        * The value can contain only Latin letters, numbers, and hyphens (`-`)
        * The value must start and end with a Latin letter or number
* `name_prefix` - (Optional) Creates a unique name beginning with the specified prefix.
    * _Value length:_ From 1 to 6 symbols
    * _Constraints:_
        * `name_prefix` cannot be specified if `name` is set
        * The value constraints are the same as for `name`

-> **Note** If `name` and `name_prefix` are not specified, Terraform will autogenerate a name with the prefix `tf-`.

* `protocol_version` - (Optional) The version of the protocol (only for HTTP).
    * _Valid values:_ `HTTP1`, `HTTP2`
* `tags` - (Optional, Editable) Map of tags to assign to the target group.
  If a provider [`default_tags` configuration block][default-tags] is used,
  tags with matching keys will overwrite those defined at the provider level.
* `target_type` - (Optional) The type of the target. All targets in a target group must have the same type.
    * _Valid values:_ `instance`
    * _Default value:_ `instance`

### health_check

The `health_check` block has the following structure:

* `enabled` - (Optional, Editable) Indicates whether health check will be enabled.
    * _Default value:_ `true`
    * _Constraints:_ Health check must be enabled for target groups with instance target type
* `healthy_threshold` - (Optional, Editable) Number of consecutive successful health checks after which the target status changes to healthy.
    * _Valid values:_ From 2 to 10
    * _Constraints:_ For TCP and UDP target groups, `healthy_threshold` and `unhealthy_threshold` must be the same
* `interval` - (Optional, Editable) The amount of time in seconds between health checks on an individual target.
    * _Valid values:_ From 5 to 300
    * _Default value:_ 30
* `port` - (Optional, Editable) The port used to perform health checks on targets.
  If value is `traffic-port`, the port, on which the target receives requests, is used for health checks.
    * _Valid values:_ `traffic-port`
    * _Default value:_ `traffic-port`
* `protocol` - (Optional, Editable) The protocol used to perform health checks on targets.
    * _Valid values:_ `TCP`, `UDP`, `HTTP`
    * _Constraints:_ `health_check.protocol` and target group `protocol` must be the same
* `timeout` - (Optional, Editable) The amount of time in seconds, after which no response indicates a failed health check.
    * _Valid values:_ From 2 to 120
    * _Constraints:_ You can set a custom timeout value for HTTP target groups only
* `unhealthy_threshold` - (Optional, Editable) Number of consecutive failed health checks after which the target status changes to unhealthy.
    * _Valid values:_ From 2 to 10
    * _Constraints:_ For TCP and UDP target groups, `healthy_threshold` and `unhealthy_threshold` must be the same

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the target group.
* `health_check` - The health check configuration.
  Exported attributes for the block are [described below](#health_check-attributes).
* `id` - The ARN of the target group.
* `tags_all` - Map of tags assigned to the target group,
  including those inherited from the provider [`default_tags` configuration block][default-tags].

#### health_check attributes

In addition to [`health_check`](#health_check) arguments above, the following attributes are exported:

* `matcher` - The HTTP code used to check the target availability.
* `path` - The destination for the health check request.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`arn_suffix`, `connection_termination`, `deregistration_delay`, `lambda_multi_value_headers_enabled`, `load_balancing_algorithm_type`, `preserve_client_ip`, `proxy_protocol_v2`, `slow_start`, `stickiness`.

## Import

The target group can be imported using `arn`, e.g.,

```
$ terraform import aws_lb_target_group.example arn:c2:elasticloadbalancing::project-name@customer-name:targetgroup/tg-12345678
```
