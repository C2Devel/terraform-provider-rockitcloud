---
subcategory: "VPC (Virtual Private Cloud)"
layout: "aws"
page_title: "aws_default_vpc"
description: |-
  Manages the default VPC.
---

[default-tags]: https://www.terraform.io/docs/providers/aws/index.html#default_tags-configuration-block

# Resource: aws_default_vpc

Manages the default VPC.

~> **Note** **This is an advanced resource** and has special caveats to be aware of when using it.
Please read this document in its entirety before using this resource.
The `aws_default_vpc` resource behaves differently from normal resources in that if a default VPC exists, Terraform does not _create_ this resource, but instead "adopts" it into management.

If no default VPC exists, Terraform creates a new default VPC, which leads to the implicit creation of other resources.
By default, `terraform destroy` does not delete the default VPC but does remove the resource from Terraform state.
Set the `force_destroy` argument to `true` to delete the default VPC.

## Example usage

### Basic example

```terraform
resource "aws_default_vpc" "default" {
  tags = {
    Name = "Default VPC"
  }
}
```

## Argument reference

The following arguments are supported:

* `enable_dns_support` - (Optional, Editable, Boolean) A flag to enable or disable the DNS support in the VPC.
    * _Default value:_ `true`
* `force_destroy` - (Optional, Editable, Boolean) Indicates whether the default VPC should be deleted during `terraform destroy`.
    * _Default value:_ `false`
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the resource.
    If the [`default_tags` configuration block][default-tags] block is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the VPC.
* `cidr_block` - (String) The IPv4 CIDR block for the VPC.
* `default_network_acl_id` - (String) The ID of the network ACL created by default on VPC creation.
* `default_route_table_id` - (String) The ID of the route table created by default on VPC creation.
* `default_security_group_id` - (String) The ID of the security group created by default on VPC creation.
* `dhcp_options_id` - (String) The ID of the DHCP options set associated to the VPC.
* `existing_default_vpc` - (Boolean) Indicates whether the default VPC already existed within the project before applying the resource.
* `id` - (String) The ID of the VPC.
* `main_route_table_id` - (String) The ID of the main route table associated with this VPC. Note that you can change a VPC's main route table by using an [`aws_main_route_table_association`](main_route_table_association.md).
* `tags_all` - (Map of strings) Key-value pairs assigned to the resource, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attributes are not currently supported:

`assign_generated_ipv6_cidr_block`, `enable_classiclink`, `enable_classiclink_dns_support`, `enable_dns_hostnames`, `instance_tenancy`, `ipv6_association_id`, `ipv6_cidr_block`, `ipv6_cidr_block_network_border_group`, `ipv6_ipam_pool_id`, `ipv6_netmask_length`, `owner_id`.

## Timeouts

Timeouts usage for the default VPCs is not currently supported.

## Import

Default VPCs can be imported using the `id`, for example:

```
$ terraform import aws_default_vpc.default vpc-12345678
```
