---
subcategory: "Route 53"
layout: "aws"
page_title: "aws_route53_zone"
description: |-
  Manages a Route53 hosted zone.
---

[default-tags]: https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block
[aws-delegation-sets]: https://docs.aws.amazon.com/Route53/latest/APIReference/actions-on-reusable-delegation-sets.html

# Resource: aws_route53_zone

Manages a Route53 hosted zone.

## Example usage

### Public zone

```terraform
resource "aws_route53_zone" "primary" {
  name = "example.com"
}
```

### Public subdomain zone

For use in subdomains, note that you need to create an `aws_route53_record` of the type `NS` as well as the subdomain zone.

```terraform
resource "aws_route53_zone" "main" {
  name = "example.com"
}

resource "aws_route53_zone" "dev" {
  name = "dev.example.com"

  tags = {
    Environment = "dev"
  }
}

resource "aws_route53_record" "dev-ns" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "dev.example.com"
  type    = "NS"
  ttl     = "30"
  records = aws_route53_zone.dev.name_servers
}
```

### Private zone

~> **Note** Each private zone must be associated with a single VPC.

```terraform
resource "aws_route53_zone" "private" {
  name = "example.com"

  vpc {
    vpc_id = aws_vpc.example.id
  }
}
```

## Argument reference

The following arguments are required:

* `name` - (Required, Forces new resource, String) The name of the hosted zone.

The following arguments are optional:

* `comment` - (Optional, Editable, String) A comment for the hosted zone.
    * _Default value:_ `Managed by Terraform`
* `force_destroy` - (Optional, Editable, Boolean) Indicates whether to destroy all records (possibly managed outside of Terraform) in the zone when the zone is destroyed.
* `tags` - (Optional, Editable, Map of strings) Key-value pairs to assign to the hosted zone.
  If the [`default_tags` configuration block][default-tags] is used within a provider configuration, the tags with matching keys will overwrite those defined at the provider level.
* `vpc` - (Optional, Editable, [Block](#vpc)) Configuration block(s) specifying a VPC to associate with a private hosted zone.

### vpc

The following arguments are required:

* `vpc_id` - (Required, Editable, String) The ID of the VPC to associate.

The following arguments are optional:

* `vpc_region` - (Optional, Editable, String) The region of the VPC to associate.

## Attribute reference

### Supported attributes

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the hosted zone.
* `name_servers` - (List of strings) A list of name servers in associated (or default) delegation set.
  Find more about delegation sets in [AWS docs][aws-delegation-sets].
* `tags_all` - (Map of strings) Key-value pairs assigned to the hosted zone, including any tags inherited from the [`default_tags` configuration block][default-tags] if used within a provider configuration.
* `zone_id` - (String) The ID of the hosted zone.
  This can be referenced by zone records.

### Unsupported attributes

~> **Note** This attribute may be present in the `terraform.tfstate` file, but it has a preset value and cannot be specified in configuration files.

The following attribute is not currently supported:

`delegation_set_id`.

## Timeouts

Timeouts usage for hosted zones is not currently supported.

## Import

Route53 hosted zone can be imported using `id`, for example:

```
$ terraform import aws_route53_zone.myzone z-xxxxxxxx
```


