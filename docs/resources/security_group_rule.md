---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_security_group_rule"
description: |-
  Manages a security group rule.
---

[iana-protocol-numbers]: https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml

# Resource: aws_security_group_rule

Manages a security group rule.
Represents a single `ingress` or `egress` group rule, which can be added to external security groups.

~> **Note** Terraform currently provides both a standalone security group rule resource (a single `ingress` or `egress` rule), and a [`aws_security_group`](security_group.md) with `ingress` and `egress` rules defined inline.
At this time you cannot use a security group with inline rules in conjunction with any security group rule resources.
Doing so will cause a conflict of rule settings and will overwrite rules.

~> **Note** Setting `protocol` to `all` or `-1` will result in the EC2 API creating a security group rule with all ports open.
This API behavior cannot be controlled by Terraform and may generate warnings in the future.

## Example usage

### Specific example

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_security_group" "example" {
  name        = "test_security_group"
  description = "test_security_group"
  vpc_id      = aws_vpc.example.id
}

resource "aws_security_group_rule" "example" {
  type              = "ingress"
  from_port         = 0
  to_port           = 65535
  protocol          = "tcp"
  cidr_blocks       = [aws_vpc.example.cidr_block]
  security_group_id = aws_security_group.example.id
}
```

## Argument reference

The following arguments are required:

* `from_port` - (Required, Forces new resource, Integer) The start of the port range (or ICMP message type number if the `protocol` value is `icmp`).
* `protocol` - (Required, Forces new resource, String) The protocol to match.
    * _Constraints:_ If the value is not `icmp`, `tcp`, `udp`, or `all`, then use the [protocol number][iana-protocol-numbers]
* `security_group_id` - (Required, Forces new resource, String) The ID of the security group to apply this rule to.
* `to_port` - (Required, Forces new resource, Integer) The end of the port range (or ICMP message code if the `protocol` value is `icmp`).
* `type` - (Required, Forces new resource, String) The type of the rule being created.
    * _Valid values:_ `egress`, `ingress`

The following arguments are optional:

* `cidr_blocks` - (Optional, Forces new resource, List of strings) The list of CIDR blocks.
    * _Constraints:_ Cannot be specified with `source_security_group_id` or `self`
* `description` - (Optional, Editable, String) The description of the rule.
* `ipv6_cidr_blocks` - (Optional, Forces new resource, List of strings) The list of IPv6 CIDR blocks.
    * _Constraints:_ Cannot be specified with `source_security_group_id` or `self`
* `self` - (Optional, Forces new resource, Boolean) Indicated whether the security group itself will be added as a source to this ingress rule.
    * _Default value:_ `false`
    * _Constraints:_ Cannot be specified with `cidr_blocks`, `ipv6_cidr_blocks`, or `source_security_group_id`
* `source_security_group_id` - (Optional, Forces new resource, String) The ID of the security group to allow access to/from, depending on the `type`.
    * _Constraints:_ Cannot be specified with `cidr_blocks`, `ipv6_cidr_blocks`, or `self`

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attribute is exported:

* `id` - (String) The ID of the security group rule.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `prefix_list_ids`.

## Timeouts

Timeouts usage for security group rules is not currently supported.

## Import

Import of the security group rules is not currently supported.
