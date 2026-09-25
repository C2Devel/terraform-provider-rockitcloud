---
subcategory: "ELB (Elastic Load Balancing)"
layout: "aws"
page_title: "aws_lb_listener"
description: |-
  Manages a listener for a load balancer.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[elb]: https://docs.k2.cloud/en/services/elb/overview.html

# Resource: aws_lb_listener

Manages a listener for a load balancer.
For details about listeners, see the [user documentation][elb].

## Example Usage

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = "tf-vpc"
  }
}

resource "aws_subnet" "example" {
  vpc_id     = aws_vpc.example.id
  cidr_block = "10.1.1.0/24"

  tags = {
    Name = "tf-subnet"
  }
}

resource "aws_lb" "example" {
  name               = "tf-lb"
  internal           = true
  load_balancer_type = "application"
  subnets            = [aws_subnet.example.id]

  tags = {
    Name = "tf-lb"
  }
}

resource "aws_lb_target_group" "example" {
  name = "tf-lb-tg"

  target_type = "instance"
  port        = 1234
  protocol    = "HTTP"
  vpc_id      = aws_vpc.example.id

  tags = {
    Name = "tf-lb-tg"
  }
}

resource "aws_lb_listener" "example" {
  load_balancer_arn = aws_lb.example.arn

  port     = 4321
  protocol = "HTTP"

  default_action {
    type = "forward"

    forward {
      target_group {
        arn = aws_lb_target_group.example.arn
      }
    }
  }

  tags = {
    Name = "tf-lb-listener"
  }
}
```

## Argument Reference

* `default_action` - (Required, Editable) The default action that will be applied to incoming requests.
  The structure of this block is [described below](#default_action).
* `load_balancer_arn` - (Required) The Amazon Resource Name (ARN) of the load balancer.
    * _ARN Format:_ `arn:c2:elasticloadbalancing::<project-name>@<customer-name>:loadbalancer/<app|net>/lb-12345678`
* `port` - (Required, Editable) The port on which the listener will receive requests.
    * _Valid values:_ From 1 to 65535
* `protocol` - (Required, Editable) The protocol for a client connection to the load balancer.
    * _Valid values:_
        * For network load balancers: `TCP`, `UDP`
        * For application load balancers: `HTTP`, `HTTPS`
* `certificate_arn` - (Optional, Editable) The ARN of the IAM server certificate.
    * _ARN Format:_ `arn:c2:iam::<customer-name>:certificate/<certificate-name>`
    * _Constraints:_ `certificate_arn` is required if `protocol` is `HTTPS`
* `tags` - (Optional, Editable) Map of tags to assign to the listener.
  If a provider [`default_tags` configuration block][default-tags] is used,
  tags with matching keys will overwrite those defined at the provider level.

### default_action

The `default_action` block has the following structure:

* `type` - (Required, Editable) The type of the routing action.
    * _Valid values:_ `forward`
* `forward` - (Optional, Editable) The block with information about forwarding requests to target groups.
  The structure of this block is [described below](#forward).
    * _Constraints:_ `forward` can be specified only if `type` is `forward` and `target_group_arn` is not specified
* `order` - (Optional, Editable) The sequential number of the action.
    * _Valid values:_ From 1 to 50000
* `target_group_arn` - (Optional, Editable, **Deprecated**) The ARN of the target group to forward traffic to.
    * _ARN Format:_ `arn:c2:elasticloadbalancing::<project-name>@<customer-name>:targetgroup/tg-12345678`
    * _Constraints:_ `target_group_arn` can be specified only if `type` is `forward` and the `forward` block is not specified

~> **Note** The argument `target_group_arn` is marked as deprecated. Use the `forward` block instead.

#### forward

The `forward` block has the following structure:

* `target_group` - (Required, Editable) List of target groups to forward traffic to.
  The structure of this block is [described below](#target_group).
    * _List size:_ From 1 to 5 elements

##### target_group

The `target_group` block has the following structure:

* `arn` - (Required, Editable) The ARN of the target group to forward traffic to.
    * _ARN Format:_ `arn:c2:elasticloadbalancing::<project-name>@<customer-name>:targetgroup/tg-12345678`
* `weight` - (Optional, Editable) The weight of the target group.
    * _Valid values:_ From 0 to 256
    * _Default value:_ 1

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the listener.
* `id` - The ARN of the listener.
* `tags_all` - Map of tags assigned to the listener,
  including those inherited from the provider [`default_tags` configuration block][default-tags].

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`alpn_policy`, `default_action.authenticate_cognito`, `default_action.authenticate_oidc`, `default_action.fixed_response`, `default_action.redirect`, `forward.stickiness`, `ssl_policy`.

## Import

The listener can be imported using `arn`, e.g.,

```
$ terraform import aws_lb_listener.example arn:c2:elasticloadbalancing::project-name@customer-name:listener/app/lb-12345678/li-12345678
```
