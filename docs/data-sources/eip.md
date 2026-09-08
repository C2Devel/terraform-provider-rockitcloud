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

## Example usage

### Search by allocation ID

```terraform
data "aws_eip" "by_allocation_id" {
  id = "eipalloc-12345678"
}
```

### Search by filters

```terraform
data "aws_eip" "by_filter" {
  filter {
    name   = "tag:Name"
    values = ["exampleNameTagValue"]
  }
}
```

### Search by public IP

```terraform
data "aws_eip" "by_public_ip" {
  public_ip = "1.2.3.4"
}
```

### Search by tags

```terraform
data "aws_eip" "by_tags" {
  tags = {
    Name = "exampleNameTagValue"
  }
}
```

## Argument reference

The arguments of this data source act as filters for querying the available Elastic IPs.
The given filters must match exactly one Elastic IP whose data will be exported as attributes.

* `filter` - (Optional, [Block](#filter)) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in the [EC2 API documentation][describe-addresses]
* `id` - (Optional, String) The ID of the allocation of the specific VPC Elastic IP to retrieve.
* `public_ip` - (Optional, String) The public IP of the specific Elastic IP to retrieve.
* `tags` - (Optional, Map of strings) Key-value pairs. Must exactly match pairs on the desired resource.

### filter

* `name` - (Required, String) The name of the filter.
    * _Constraints:_ Filter names are case-sensitive
* `values` - (Required, List of strings) One or more filter values.
    * _Constraints:_ Filter values are case-sensitive

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `association_id` - (String) The ID of the address association with an instance in a VPC.
* `domain` - (String) Indicates whether the address is for use in EC2 Classic (standard) or in a VPC (vpc).
* `id` - (String) If VPC Elastic IP, the allocation identifier.
* `instance_id` - (String) The ID of the instance that the address is associated with (if any).
* `network_interface_id` - (String) The ID of the network interface.
* `network_interface_owner_id` - (String) The ID of the project the network interface belongs to.
* `private_ip` - (String) The private IP address associated with the Elastic IP address.
* `public_ip` - (String) The public IP address of the Elastic IP address.
* `public_ipv4_pool` - (String) The ID of an address pool.
* `tags` - (Map of strings) Key-value pairs assigned to the Elastic IP address.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`carrier_ip`, `customer_owned_ip`, `customer_owned_ipv4_pool`, `private_dns`, `public_dns`.

~> **Note** The data source computes the `public_dns` and `private_dns` attributes according to the [AWS VPC DNS Guide][vpc-dns-hostnames] as they are not available with the EC2 API.
