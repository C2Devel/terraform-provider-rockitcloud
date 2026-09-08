---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_eip"
description: |-
  Manages an Elastic IP.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[elastic-ips]: https://docs.k2.cloud/en/services/networking/addresses/operations.html
[timeouts]: https://developer.hashicorp.com/terraform/plugin/framework/resources/timeouts
[vpc-dns-guide]: https://docs.aws.amazon.com/vpc/latest/userguide/vpc-dns.html#vpc-dns-hostnames

# Resource: aws_eip

Manages an Elastic IP address. For more information about EIPs, see the [user documentation][elastic-ips].

## Example usage

### Single EIP associated with an instance

```terraform
resource "aws_eip" "example" {
  instance = "i-12345678"
  vpc      = true
}
```

### Attaching an EIP to an instance with a pre-assigned private IP

```terraform
resource "aws_vpc" "default" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "tf_test_subnet" {
  vpc_id     = aws_vpc.default.id
  cidr_block = "10.0.0.0/24"
}

resource "aws_instance" "foo" {
  ami           = "cmi-12345678" # add image id, change instance type if needed
  instance_type = "m1.micro"

  private_ip = "10.0.0.12"
  subnet_id  = aws_subnet.tf_test_subnet.id
}

resource "aws_eip" "bar" {
  vpc = true

  associate_with_private_ip = "10.0.0.12"
}
```

### Allocating the EIP address from the BYOIP pool

```terraform
resource "aws_eip" "byoip-ip" {
  vpc              = true
  public_ipv4_pool = "ipv4pool-ec2-012345"
}
```

## Argument reference

* `address` - (Optional, Forces new resource, String) An IP address from an EC2 BYOIP pool.
    * _Constraints:_ This option is only available for EIPs in a VPC
* `associate_with_private_ip` - (Optional, Editable, String) A user-specified primary or secondary private IP address to associate with the Elastic IP address.
    * _Constraints:_ If no private IP address is specified, the Elastic IP address is associated with the primary private IP address
* `instance` - (Optional, Editable, String) The ID of the EC2 instance.
* `network_interface` - (Optional, Editable, String) The ID of the network interface to associate with.
* `public_ipv4_pool` - (Optional, Forces new resource, String) The ID of the EC2 IPv4 address pool.
    * _Constraints:_ This option is only available for EIPs in a VPC
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the Elastic IP address. If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
    * _Constraints:_ Tags can only be applied to EIPs in a VPC
* `vpc` - (Optional, Forces new resource, Boolean) Indicates whether the EIP is in a VPC.

~> **Note** You can specify either the ID of `instance` or the ID of `network_interface`, but not both.

~> **Note** If both `public_ipv4_pool` and `address` are specified, `address` will be used in the case both options are defined as API only requires one or the other.

## Attribute reference

### Supported attributes

~> **Note** The data source computes the `public_dns` and `private_dns` attributes according to the [AWS VPC DNS Guide][vpc-dns-guide] as they are not available with the EC2 API.

In addition to all arguments above, the following attributes are exported:

* `allocation_id` - (String) The ID representing the allocation of the IP address.
* `association_id` - (String) The ID representing the association of the allocation of the IP address with an instance or a private IP address.
* `domain` - (String) Indicates if this EIP is for use in VPC (`vpc`).
* `id` - (String) The ID of the EIP allocation.
* `private_ip` - (String) The private IP address.
  Can be `""` if `associate_with_private_ip` is specified.
* `public_ip` - (String) The public IP address.
* `tags_all` - (Map of strings) Key-value pairs assigned to the Elastic IP, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`carrier_ip`, `customer_owned_ip`, `customer_owned_ipv4_pool`, `network_border_group`, `private_dns`, `public_dns`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `read` - (Default `15 minutes`) Used when querying for information about EIPs.
* `update` - (Default `5 minutes`) Used when updating an EIP.
* `delete` - (Default `3 minutes`) Used when deleting an EIP.

## Import

EIPs in a VPC can be imported using their allocation ID, for example:

```
$ terraform import aws_eip.bar eipalloc-1234567
```

EIPs can be imported using their public IP, for example:

```
$ terraform import aws_eip.bar 1.1.1.1
```
