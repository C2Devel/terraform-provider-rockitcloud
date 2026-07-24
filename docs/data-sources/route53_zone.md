---
subcategory: "Route 53"
layout: "aws"
page_title: "aws_route53_zone"
description: |-
  Provides information about a Route 53 hosted zone.
---

# Data Source: aws_route53_zone

Provides information about a Route 53 hosted zone.

This data source allows you to find a hosted zone ID given hosted zone name and certain search criteria.

## Example usage

The following example shows how to get a hosted zone from its name and how to create a record set from this data.


```terraform
data "aws_route53_zone" "selected" {
  name         = "test.com."
  private_zone = true
}

resource "aws_route53_record" "www" {
  zone_id = data.aws_route53_zone.selected.zone_id
  name    = "www.${data.aws_route53_zone.selected.name}"
  type    = "A"
  ttl     = "300"
  records = ["10.0.0.1"]
}
```

## Argument reference

The arguments of this data source act as filters for querying the available hosted zones.
The given filter must match exactly one hosted zone.

~> **Note** You have to specify either `zone_id` or `name`, but not both of them.
If you use the `name` field to search for a private hosted zone, you need to set the argument `private_zone` value to `true`.

* `name` - (Optional, String) The name of the desired hosted zone.
* `private_zone` - (Optional, Boolean) Indicates whether the hosted zone is private.
  Used with the `name` field to get a private hosted zone.
    * _Default value:_ `false`
* `resource_record_set_count` - (Optional, Integer) The number of record sets in the hosted zone.
  Used with `name` field.
* `vpc_id` - (Optional, String) Used with `name` field to get a private hosted zone associated with the `vpc_id`.
  In this case, `private_zone` is not required.
* `zone_id` - (Optional, String) The ID of the desired hosted zone.

## Attribute reference

### Supported attributes

This data source will complete the data by populating any fields that are not included in the configuration with the data for the selected hosted zone.

In addition to all arguments above, the following attributes are exported:

* `arn` - (String) The Amazon Resource Name (ARN) of the hosted zone.
* `caller_reference` - (String) Caller Reference of the hosted zone.
* `comment` - (String) The comment field of the hosted zone.
* `name_servers` - (List of strings) The list of DNS name servers for the hosted zone.
* `tags` - (Map of strings) Key-value pairs assigned to the hosted zone.

### Unsupported attributes

~> **Note** These attributes may be present in the `terraform.tfstate` file, but they have preset values and cannot be specified in configuration files.

The following attributes are not currently supported:

`linked_service_description`, `linked_service_principal`.
