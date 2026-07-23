---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_network_acl"
description: |-
  Creates a network ACL.
---

[attribute-as-blocks]: https://www.terraform.io/docs/configuration/attr-as-blocks.html
[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[icmp-parameters]:  https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml
[network-acl]: https://docs.k2.cloud/en/services/security/networkacl.html

# Resource: aws_network_acl

Creates a network ACL.
You might set up network ACLs with rules similar to your security groups to add another layer of security to your VPC.

~> **Note on network ACLs and network ACL rules** Terraform currently provides both a standalone [`aws_network_acl_rule`](network_acl_rule.md) resource and a network ACL resource with rules defined inline.
At this time you cannot use a network ACL with inline rules in conjunction with any network ACL rule resources.
Doing so will cause a conflict of rule settings and will overwrite rules.

~> **Note on network ACLs and network ACL associations** Terraform provides both a standalone [aws_network_acl_association](network_acl_association.md) resource and a network ACL resource with a `subnet_ids` attribute.
Do not use the same subnet ID in both a network ACL resource and a network ACL association resource.
Doing so will cause a conflict of associations and will overwrite the association.

For more information about network ACLs, see the documentation on [Network ACL][network-acl].

## Example usage

### Specific example

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_network_acl" "example" {
  vpc_id = aws_vpc.example.id

  egress {
    protocol   = "tcp"
    rule_no    = 200
    action     = "allow"
    cidr_block = "10.3.0.0/18"
    from_port  = 443
    to_port    = 443
  }

  ingress {
    protocol   = "tcp"
    rule_no    = 100
    action     = "allow"
    cidr_block = "10.3.0.0/18"
    from_port  = 80
    to_port    = 80
  }

  tags = {
    Name = "main"
  }
}
```

## Argument reference

The following arguments are required:

* `vpc_id` - (Required, Forces new resource, String) The ID of the associated VPC.

The following arguments are optional:

* `egress` - (Optional, Editable, [Block](#egress-and-ingress)) One or more egress rules (for outgoing traffic).
* `ingress` - (Optional, Editable, [Block](#egress-and-ingress)) One or more ingress rules (for incoming traffic).
* `subnet_ids` - (Optional, Editable, List of strings) The list of subnet IDs to apply the ACL to.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource.
    If the [`default_tags` configuration block][default-tags] block is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

### egress and ingress

Both arguments are processed in [attribute-as-blocks mode][attribute-as-blocks].

Both `egress` and `ingress` support the following keys:

* `action` - (Required, Editable, String) The action to take.
    * _Valid values:_ `allow`, `deny`
* `from_port` - (Required, Editable, Integer) The start of the port range.
* `protocol` - (Required, Editable, String) The protocol to match.
    If using the `-1` (`all`) value, then you must specify a start and end numbers of `0`.
* `rule_no` - (Required, Editable, Integer) The rule number.
    Used for ordering.
    Rules are processed in ascending order.
    * _Valid values:_ From 1 to 32766
* `to_port` - (Required, Editable, Integer) The end of the port range.
* `cidr_block` - (Optional, Editable, String) The CIDR block to match. This must be a valid network mask.
* `icmp_code` - (Optional, Editable, Integer) The ICMP message code to be used.
    * _Default value:_ `0`
* `icmp_type` - (Optional, Editable, Integer) The ICMP message type to be used.
    * _Default value:_ `0`

~> **Note** For more information on ICMP message types and codes, see [ICMP Parameters][icmp-parameters]

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the network ACL.
* `id` - (String) The ID of the network ACL.
* `tags_all` - (Optional, Map of strings) Key-value pairs assigned to the resource,  including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`egress.ipv6_cidr_block`, `ingress.ipv6_cidr_block`, `owner_id`.

## Timeouts

Timeouts usage for the network ACLs is not currently supported.

## Import

Network ACLs can be imported using `id`, for example:

```
$ terraform import aws_network_acl.main acl-12345678
```

