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

# Resource: aws_eip

Manages an Elastic IP (EIP). For more information about EIPs, see [user documentation][elastic-ips].

## Example usage

### EIP associated with an instance

```terraform
resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "example" {
  vpc_id     = aws_vpc.example.id
  cidr_block = "10.0.0.0/24"
}

resource "aws_internet_gateway" "example" {
  vpc_id = aws_vpc.example.id
}

resource "aws_route" "default_route" {
  route_table_id         = aws_vpc.example.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.example.id
}

resource "aws_instance" "example" {
  ami           = "cmi-12345678"
  instance_type = "m1.micro"
  subnet_id     = aws_subnet.example.id
}

resource "aws_eip" "example" {
  depends_on = [aws_internet_gateway.example]

  instance = aws_instance.example.id

  tags = {
    Name = "tf-eip"
  }
}
```

### EIP associated with a network interface

~> **Note** This example uses the same VPC, subnet and internet gateway as in the [EIP associated with an instance example](#eip-associated-with-an-instance).

```terraform
resource "aws_network_interface" "example" {
  subnet_id = aws_subnet.example.id
}

resource "aws_eip" "example" {
  depends_on = [aws_internet_gateway.example]

  network_interface = aws_network_interface.example.id

  tags = {
    Name = "tf-eip"
  }
}
```

### Allocating an EIP from the BYOIP pool

```terraform
resource "aws_eip" "example" {
  public_ipv4_pool = "ipv4pool-ec2-012345"
}
```

## Argument reference

The following arguments are optional:

* `address` - (Optional, Forces new resource, String) An IP address from an EC2 BYOIP pool.
* `instance` - (Optional, Editable, String) The ID of the EC2 instance. The EIP is associated with the primary network interface of the instance.
    * _Constraints:_ Conflicts with the `network_interface` argument
* `network_interface` - (Optional, Editable, String) The ID of the network interface to associate with.
    * _Constraints:_ Conflicts with the `instance` argument
* `public_ipv4_pool` - (Optional, Forces new resource, String) The ID of the EC2 IPv4 address pool.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the EIP. If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
* `vpc` - (Optional, Editable, Boolean, **Deprecated**) Indicates whether the EIP is in a VPC.

~> **Note** The argument `vpc` is deprecated.
Its value is ignored: all EIPs are for use in a VPC.

~> **Note** If both `public_ipv4_pool` and `address` are specified, `address` will be used in the case both options are defined as API only requires one or the other.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `allocation_id` - (String) The ID of the allocation of the IP address.
* `association_id` - (String) The ID of the address association with an instance.
    * _Constraints:_ Empty unless the EIP is attached to an instance or a network interface
* `domain` - (String) The domain in which the EIP is used.
    * _Constraints:_ Always `vpc`
* `id` - (String) The ID of the EIP allocation.
* `private_dns` - (String) The private DNS name of the network interface this EIP is attached to.
    * _Constraints:_ Empty unless the EIP is attached to an instance or a network interface
* `private_ip` - (String) The private IP address.
    * _Constraints:_ Empty unless the EIP is attached to an instance or a network interface
* `public_dns` - (String) The public DNS name of the network interface this EIP is attached to.
    * _Constraints:_ Empty unless the EIP is attached to an instance or a network interface
* `public_ip` - (String) The public IP address.
* `tags_all` - (Map of strings) Key-value pairs assigned to the EIP, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`associate_with_private_ip`, `carrier_ip`, `customer_owned_ip`, `customer_owned_ipv4_pool`, `network_border_group`.

## Timeouts

The `timeouts` block allows you to specify [timeouts] for certain actions:

* `read` - (Default `15 minutes`) Used when waiting for a newly created EIP to become visible.
* `update` - (Default `5 minutes`) Used when associating the EIP with an instance or a network interface.
* `delete` - (Default `3 minutes`) Used when releasing an EIP.

## Import

The EIP can be imported using the allocation ID, for example:

```
$ terraform import aws_eip.example eipalloc-1234567
```

The EIP can also be imported using the public IP, for example:

```
$ terraform import aws_eip.example 1.1.1.1
```
