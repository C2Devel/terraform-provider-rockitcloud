---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_eip_association"
description: |-
  Manages an EIP association.
---

# Resource: aws_eip_association

Manages an EIP association as a top level resource, to associate and
disassociate Elastic IPs from instances and network interfaces.

~> **Note** `aws_eip_association` is useful in scenarios where EIPs are either
pre-existing or distributed to customers or users and therefore cannot be changed.

## Example usage

```terraform
resource "aws_eip_association" "eip_assoc" {
  instance_id   = aws_instance.web.id
  allocation_id = aws_eip.example.id
}

resource "aws_instance" "web" {
  ami               = "cmi-12345678" # add image id, change instance type if needed
  availability_zone = "ru-msk-vol52"
  instance_type     = "m1.micro"

  tags = {
    Name = "HelloWorld"
  }
}

resource "aws_eip" "example" {
  vpc = true
}
```

## Argument reference

The following arguments are supported:

* `allocation_id` - (Optional, Forces new resource, String) The ID of the allocation.
    * _Constraints:_ Required if the `public_ip` is not supplied
* `allow_reassociation` - (Optional, Forces new resource, Boolean) Indicates whether to allow an Elastic IP to be re-associated.
    * _Default value:_ `true`
* `instance_id` - (Optional, Forces new resource, String) The ID of the instance.
    * _Constraints:_ Required if the `network_interface_id` is not supplied
* `network_interface_id` - (Optional, Forces new resource, String) The ID of the network interface.
    * _Constraints:_ Required if the `instance_id` is not supplied
* `public_ip` - (Optional, Forces new resource, String) The Elastic IP address.
    * _Constraints:_ Required if the `allocation_id` is not supplied

## Attribute reference

In addition to all arguments above, the following attributes are exported:

* `association_id` - (String) The ID that represents the association of the Elastic IP address with an instance.
* `allocation_id` - (String) The ID of the allocation.
* `instance_id` - (String) The ID of the instance that the address is associated with.
* `network_interface_id` - (String) The ID of the network interface.
* `private_ip_address` - (String) The private IP address associated with the Elastic IP address.
* `public_ip` - (String) The public IP address of the Elastic IP.

## Timeouts

Timeouts usage for EIP association is not currently supported.

## Import

EIP associations can be imported using IDs of their associations, for example:

```
$ terraform import aws_eip_association.test eipassoc-12345678
```
