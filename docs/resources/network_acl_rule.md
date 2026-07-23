---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_acl_rule"
description: |-
  Creates a network ACL rule.
---

[icmp-parameters]: https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml
[protocol-numbers]: https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml

# Resource: aws_network_acl_rule

Creates an entry, or a rule, in a network ACL with the specified rule number.

~> **Note on network ACLs and network ACL rules**
     Terraform currently provides both a standalone network ACL rule resource and an [aws_network_acl](network_acl.md) resource with rules defined inline.
     At this time you cannot use a network ACL with inline rules in conjunction with any network ACL rule resources.
     Doing so will cause a conflict of rule settings and will overwrite rules.

## Example usage

### Specific example

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_network_acl" "example" {
  vpc_id = aws_vpc.example.id
}

resource "aws_network_acl_rule" "example" {
  network_acl_id = aws_network_acl.example.id
  rule_number    = 200
  egress         = false
  protocol       = "tcp"
  rule_action    = "allow"
  cidr_block     = aws_vpc.example.cidr_block
  from_port      = 22
  to_port        = 22
}
```

## Argument reference

The following arguments are required:

* `cidr_block` - (Required, Forces new resource, String) The network range to allow or deny, in CIDR notation.
    * _Example:_ `172.16.0.0/24`
* `network_acl_id` - (Required, Forces new resource, String) The ID of the network ACL.
* `protocol` - (Required, Forces new resource, String) The protocol. A value of `-1` means all protocols.
* `rule_action` - (Required, Forces new resource, String) Indicates whether to allow or deny the traffic that matches the rule.
    * _Valid values:_ `allow`, `deny`
* `rule_number` - (Required, Forces new resource, Integer) The rule number.
    Used for ordering.
    ACL rules are processed in ascending order.
    * _Valid values:_ From 0 to 32766

The following arguments are optional:

* `egress` - (Optional, Forces new resource, Boolean) Indicates whether this is an egress rule (rule is applied to outgoing traffic).
    * _Default value:_ `false`
* `from_port` - (Optional, Forces new resource, Integer) The start of the port range.
* `icmp_type` - (Optional, Forces new resource, Integer) The ICMP message type.
    * _Example:_ `-1`
    * _Constraints:_ Required if specifying ICMP for the protocol
* `icmp_code` - (Optional, Forces new resource, Integer) The ICMP message code.
    * _Example:_ `-1`
    * _Constraints:_ Required if specifying ICMP for the protocol
* `to_port` - (Optional, Forces new resource, Integer) The end of the port range.

~> **Note** If the value of `protocol` is `-1` or `all`, the `from_port` and `to_port` values will be ignored and the rule will apply to all ports.

~> **Note** If the value of `icmp_type` is `-1` (which results in a wildcard ICMP message type), the `icmp_code` must also be set to `-1` (wildcard ICMP message code).

-> **Info** For more information on ICMP message types and codes, see [ICMP Parameters][icmp-parameters].

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attribute is exported:

* `id` - (String) The ID of the network ACL rule.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported: `ipv6_cidr_block`.

## Timeouts

Timeouts usage for the network ACL rules is not currently supported.

## Import

Individual rules can be imported using `NETWORK_ACL_ID:RULE_NUMBER:PROTOCOL:EGRESS`, where `PROTOCOL` can be a decimal (for example, `6`) or string (for example, `tcp`) value.
If importing a rule previously provisioned by Terraform, the `PROTOCOL` must be the input value used at creation time.
For more information about protocol numbers and keywords, see [Protocol numbers][protocol-numbers].

For example, import a network ACL rule with an argument like this:

```console
$ terraform import aws_network_acl_rule.my_rule acl-12345678:100:tcp:false
```

Or by the protocol's decimal value:

```console
$ terraform import aws_network_acl_rule.my_rule acl-12345678:100:6:false
```
