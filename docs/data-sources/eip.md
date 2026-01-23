---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_eip"
description: |-
  Provides information about an Elastic IP.
---

[describe-addresses]: https://docs.k2.cloud/en/api/ec2/actions/addresses/DescribeAddresses.html
[vpc-dns-hostnames]: https://docs.aws.amazon.com/vpc/latest/userguide/vpc-dns.html#vpc-dns-hostnames

# Data Source: aws_eip

Provides information about an Elastic IP.

## Example Usage

### Search By Allocation ID

```terraform
data "aws_eip" "by_allocation_id" {
  id = "eipalloc-12345678"
}
```

### Search By Filters

```terraform
data "aws_eip" "by_filter" {
  filter {
    name   = "tag:Name"
    values = ["exampleNameTagValue"]
  }
}
```

### Search By Public IP

```terraform
data "aws_eip" "by_public_ip" {
  public_ip = "1.2.3.4"
}
```

### Search By Tags

```terraform
data "aws_eip" "by_tags" {
  tags = {
    Name = "exampleNameTagValue"
  }
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available Elastic IPs.
The given filters must match exactly one Elastic IP whose data will be exported as attributes.

* `filter` - (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-addresses]
* `id` - (Optional) The ID of the allocation of the specific VPC Elastic IP to retrieve.
* `public_ip` - (Optional) The public IP of the specific Elastic IP to retrieve.
* `tags` - (Optional) Map of tags, each pair of which must exactly match a pair on the desired Elastic IP.

## Attribute Reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `association_id` - The ID representing the association of the address with an instance in a VPC.
* `domain` - Indicates whether the address is for use in EC2-Classic (standard) or in a VPC (vpc).
* `id` - If VPC Elastic IP, the allocation identifier.
* `instance_id` - The ID of the instance that the address is associated with (if any).
* `network_interface_id` - The ID of the network interface.
* `network_interface_owner_id` - The ID of the project that owns the network interface.
* `private_ip` - The private IP address associated with the Elastic IP address.
* `public_ip` - Public IP address of Elastic IP.
* `public_ipv4_pool` - The ID of an address pool.
* `tags` - Map of tags assigned to the Elastic IP.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`carrier_ip`, `customer_owned_ip`, `customer_owned_ipv4_pool`, `private_dns`, `public_dns`.

~> **Note** The data source computes the `public_dns` and `private_dns` attributes according to the [AWS VPC DNS Guide][vpc-dns-hostnames] as they are not available with the EC2 API.
