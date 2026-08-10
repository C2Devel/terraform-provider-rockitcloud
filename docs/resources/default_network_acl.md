---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_default_network_acl"
description: |-
  Manages the default network ACL of a VPC.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block
[icmp-parameters]: https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml
[network-acl]: https://docs.k2.cloud/en/services/security/networkacl.html
[protocol-number]: https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml

# Resource: aws_default_network_acl

Manages the default network ACL of a VPC.
This resource can manage the default network ACL of the default or a non-default VPC.

~> **Note** This is an advanced resource with special caveats.
Please read this document in its entirety before using this resource.
The `aws_default_network_acl` resource behaves differently from normal resources.
Terraform does not _create_ this resource but instead attempts to "adopt" it into management.

Every VPC has a default network ACL that can be managed but not destroyed.
When Terraform first adopts the default network ACL, it **immediately removes all rules in the ACL**.
It then proceeds to create any rules specified in the configuration.
This step is required so that only the rules specified in the configuration are created.

This resource treats its inline rules as absolute; only the rules defined inline are created, and any additions/removals external to this resource will result in diffs being shown.
For these reasons, this resource is incompatible with the `aws_network_acl_rule` resource.

For more information about network ACLs, see the documentation on [Network ACL][network-acl].

## Example usage

### Basic example

The following configuration gives the default network ACL the same rules that the cloud includes but pulls the resource under management by Terraform.
This means that any ACL rules added or changed will be detected as drift.

```terraform
resource "aws_vpc" "mainvpc" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_default_network_acl" "default" {
  default_network_acl_id = aws_vpc.mainvpc.default_network_acl_id

  ingress {
    protocol   = -1
    rule_no    = 100
    action     = "allow"
    cidr_block = "0.0.0.0/0"
    from_port  = 0
    to_port    = 0
  }

  egress {
    protocol   = -1
    rule_no    = 100
    action     = "allow"
    cidr_block = "0.0.0.0/0"
    from_port  = 0
    to_port    = 0
  }
}
```

### Specific example: deny all egress traffic, allow ingress traffic

The following configuration denies all outgoing traffic by omitting any `egress` rules, while including the default `ingress` rule to allow all incoming traffic.

```terraform
resource "aws_vpc" "mainvpc" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_default_network_acl" "default" {
  default_network_acl_id = aws_vpc.mainvpc.default_network_acl_id

  ingress {
    protocol   = -1
    rule_no    = 100
    action     = "allow"
    cidr_block = aws_vpc.mainvpc.cidr_block
    from_port  = 0
    to_port    = 0
  }
}
```

### Specific example: deny all traffic to any subnet in the default network ACL

This configuration denies all traffic in the default network ACL.
This can be useful if you want to lock down the VPC to force all resources to assign a non-default ACL.

```terraform
resource "aws_vpc" "mainvpc" {
  cidr_block = "10.1.0.0/16"
}

resource "aws_default_network_acl" "default" {
  default_network_acl_id = aws_vpc.mainvpc.default_network_acl_id

  # no rules defined, deny all traffic in this ACL
}
```

### Managing subnets in a default network ACL

Within a VPC, all subnets must be associated with a network ACL.
In order to "delete" the association between a subnet and a non-default network ACL, the association is destroyed by replacing it with an association between the subnet and the default ACL instead.

When managing the default network ACL, you cannot "remove" subnets.
Instead, they must be reassigned to another network ACL, or the subnet itself must be destroyed.
Because of these requirements, removing the `subnet_ids` attribute from the configuration of the `aws_default_network_acl` resource may result in a reoccurring plan, until the subnets are reassigned to another network ACL or are destroyed.

Because subnets are by default associated with the default network ACL, any non-explicit association will show up as a plan to remove the subnet.
For example, if you have a custom `aws_network_acl` with two subnets attached, and you remove the `aws_network_acl` resource, after successfully destroying this resource future plans will show a diff on the managed `aws_default_network_acl`, as those two subnets have been orphaned by the now destroyed network acl and thus adopted by the default network ACL.
To avoid a reoccurring plan, they will need to be reassigned, destroyed, or added to the `subnet_ids` attribute of the `aws_default_network_acl` entry.

As an alternative to the above, you can also specify the following lifecycle configuration in your `aws_default_network_acl` resource:

```terraform
resource "aws_default_network_acl" "default" {
  # ... other configuration ...

  lifecycle {
    ignore_changes = [subnet_ids]
  }
}
```

## Removing `aws_default_network_acl` from your configuration

Each VPC comes with a default network ACL that cannot be deleted.
The `aws_default_network_acl` allows you to manage this network ACL, but Terraform cannot destroy it.
Removing this resource from your configuration will remove it from your statefile and management, **but will not destroy the network ACL.**
All subnets associations and ingress or egress rules will be left as they are at the time of removal.
You can resume managing them via the cloud console.

## Argument reference

The following arguments are required:

* `default_network_acl_id` - (Required, Forces new resource, String) The ID of the network ACL to manage.
  This attribute is exported from `aws_vpc`, or manually found via the cloud console.

The following arguments are optional:

* `egress` - (Optional, Editable, [Block](#egress-and-ingress)) One or more egress rules (for outgoing traffic).
* `ingress` - (Optional, Editable, [Block](#egress-and-ingress)) One or more ingress rules (for incoming traffic).
* `subnet_ids` - (Optional, Editable, List of strings) The list of subnet IDs to apply the ACL to.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource.
    If the [`default_tags` configuration block][default-tags] block is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

### egress and ingress

Both `egress` and `ingress` objects have the same arguments.

The following arguments are required:

* `action` - (Required, Editable, String) The action to take.
    * _Valid values:_ `allow`, `deny`
* `from_port` - (Required, Editable, Integer) The start of the port range to which the rule applies.
* `protocol` - (Required, Editable, String) The protocol to match.
    * _Constraints:_
        * If using the `-1` value (semantically equivalent to `all`, which is not a valid value here), you must specify the `from_port` and `to_port` arguments values equal to `0`
        * If the `protocol` value is not `icmp`, `tcp`, `udp`, or `-1`, then refer to the [protocol number][protocol-number] for detailed information
* `rule_no` - (Required, Editable, Integer) The rule number. Used for ordering.
* `to_port` - (Required, Editable, Integer) The end of the port range to which the rule applies.

The following arguments are optional:

* `cidr_block` - (Optional, Editable, String) The CIDR block to match. This must be a valid network mask.
* `icmp_code` - (Optional, Editable, Integer) The ICMP message code to be used.
    * _Default value:_ `0`
* `icmp_type` - (Optional, Editable, Integer) The ICMP message type to be used.
    * _Default value:_ `0`

-> For more information on ICMP types and codes, see [Internet Control Message Protocol (ICMP) Parameters][icmp-parameters].

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the default network ACL.
* `id` - (String) The ID of the default network ACL.
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the  [`default_tags` configuration block][default-tags] if used within a provider configuration.
* `vpc_id` -  (String) The ID of the associated VPC.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`egress.ipv6_cidr_block`, `ingress.ipv6_cidr_block`, `owner_id`.

## Timeouts

Timeouts usage for the default network ACLs is not currently supported.

## Import

Default network ACLs can be imported using `id`, for example:

```
$ terraform import aws_default_network_acl.sample acl-12345678
```

