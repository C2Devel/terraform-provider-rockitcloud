---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "aws_key_pair"
description: |-
  Provides information about a key pair.
---

[describe-key-pairs]: https://docs.k2.cloud/en/api/ec2/actions/key_pairs/DescribeKeyPairs.html

# Data Source: aws_key_pair

Provides information about a key pair.

## Example Usage

The following example shows how to get a key pair from its name.

```terraform
data "aws_key_pair" "example" {
  key_name = "test"
  filter {
    name   = "tag:Component"
    values = ["web"]
  }
}

output "fingerprint" {
  value = data.aws_key_pair.example.fingerprint
}

output "name" {
  value = data.aws_key_pair.example.key_name
}

output "id" {
  value = data.aws_key_pair.example.id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available key pairs.
The given filters must match exactly one key pair whose data will be exported as attributes.

* `filter` -  (Optional) One or more name/value pairs to use as filters.
    * _Valid values:_ See supported names and values in [EC2 API documentation][describe-key-pairs]
* `key_name` - (Optional) The name of the key pair.
* `key_pair_id` - (Optional) The ID of the key pair.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the key pair.
* `fingerprint` - The SHA-1 digest of the DER encoded private key.
* `id` - The ID of the key pair.
* `tags` - Map of tags assigned to the key pair.
