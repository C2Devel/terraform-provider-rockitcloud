---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_eip_association"
description: |-
  Manages an EIP association.
---

# Resource: aws_eip_association

Manages an Elastic IP (EIP) association as a top level resource, to associate and
disassociate EIPs from instances and network interfaces.

~> **Note** `aws_eip_association` is useful in scenarios where EIPs are either
pre-existing or distributed to customers or users and therefore cannot be changed.

## Example usage

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

resource "aws_eip" "example" {}

resource "aws_eip_association" "example" {
  depends_on = [aws_internet_gateway.example]

  instance_id   = aws_instance.example.id
  allocation_id = aws_eip.example.id
}
```

## Argument reference

The following arguments are supported:

* `allocation_id` - (Optional, Forces new resource, String) The ID of the allocation.
    * _Constraints:_ Required if the `public_ip` is not supplied
* `allow_reassociation` - (Optional, Forces new resource, Boolean) Indicates whether to allow an EIP to be re-associated.
    * _Default value:_ `true`
* `instance_id` - (Optional, Forces new resource, String) The ID of the instance.
    * _Constraints:_ Required if the `network_interface_id` is not supplied
* `network_interface_id` - (Optional, Forces new resource, String) The ID of the network interface.
    * _Constraints:_ Required if the `instance_id` is not supplied
* `public_ip` - (Optional, Forces new resource, String) The EIP address.
    * _Constraints:_ Required if the `allocation_id` is not supplied

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `allocation_id` - (String) The ID of the allocation.
* `id` - (String) The ID of the association.
* `instance_id` - (String) The ID of the instance that the address is associated with.
* `network_interface_id` - (String) The ID of the network interface.
* `private_ip_address` - (String) The private IP address associated with the EIP.
* `public_ip` - (String) The public IP address of the EIP.

~> **Note** The `private_ip_address` argument was removed and cannot be specified in configuration files.
The value is set by the platform when the EIP is associated.

## Timeouts

Timeouts usage for EIP association is not currently supported.

## Import

EIP associations can be imported using IDs of their associations, for example:

```
$ terraform import aws_eip_association.example eipassoc-12345678
```
